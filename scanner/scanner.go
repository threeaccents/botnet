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

// ScanHosts is
func ScanHosts(hosts []string) []string {
	t1 := time.Now()
	wg := new(sync.WaitGroup)
	var mutex = &sync.Mutex{}
	var results []string
	for _, host := range hosts {
		wg.Add(1)
		go func(host string) {
			res := scanCommonPorts(host, wg)
			if len(res) > 0 {
				for _, r := range res {
					mutex.Lock()
					results = append(results, r)
					mutex.Unlock()
				}
			}
		}(host)
	}
	wg.Wait()

	fmt.Println("[*] done scanning scan time", time.Since(t1))

	return results
}

// scanAllPorts is
func scanAllPorts(host string) {
	for port := 0; port < 65535; port++ {
		if err := scanPort(host, strconv.Itoa(port)); err != nil {
			continue
		}

		fmt.Printf("%s:%d", host, port)
	}
}

// scanCommonPorts is
func scanCommonPorts(host string, wg *sync.WaitGroup) []string {
	defer wg.Done()
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
