package scanner

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	commonPorts = []string{"22"}
)

// Scanner is
type Scanner struct {
	wg    *sync.WaitGroup
	mutex *sync.Mutex
}

// ScanHosts is
func (s *Scanner) ScanHosts(hosts []string) <-chan string {
	if s.wg == nil {
		s.wg = new(sync.WaitGroup)
	}

	fmt.Println("[*] starting hosts scan")
	results := make(chan string)

	go func() {
		for _, host := range hosts {
			s.wg.Add(1)
			go func(host string) {
				res := s.ScanCommonPorts(host)
				if len(res) > 0 {
					for _, r := range res {
						results <- r
					}
				}
			}(host)
		}
		s.wg.Wait()
		close(results)
	}()

	return results
}

// ScanAllPorts is
func (s *Scanner) ScanAllPorts(host string) {
	for port := 0; port < 65535; port++ {
		if err := scanPort(host, strconv.Itoa(port)); err != nil {
			continue
		}

		fmt.Printf("%s:%d", host, port)
	}
}

// ScanCommonPorts is
func (s *Scanner) ScanCommonPorts(host string) []string {
	defer s.wg.Done()
	var result []string
	for _, port := range commonPorts {
		if err := scanPort(host, port); err != nil {
			continue
		}

		fmt.Printf("[*] FOUND %s:%s\n", host, port)
		result = append(result, fmt.Sprintf("%s:%s", host, port))

	}

	return result
}

func scanPort(host, port string) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 10*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}
