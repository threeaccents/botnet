package scanner

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	commonPorts = []string{"20", "21", "22", "23", "25", "53", "67", "68", "69", "80", "110", "123", "137", "138", "139", "143", "161", "162", "179", "389", "443", "636", "989", "990"}
)

// ScanHosts is
func ScanHosts(hosts []string) []string {
	wg := new(sync.WaitGroup)
	var mutex = &sync.Mutex{}
	var results []string
	for _, host := range hosts {
		go func(host string) {
			wg.Add(1)
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

	return results
}

// scanAllPorts is
func scanAllPorts(host string) {
	for port := 0; port < 65535; port++ {
		ok, err := scanPort(host, strconv.Itoa(port))
		if err != nil {
			continue
		}

		if ok {
			fmt.Printf("%s:%d", host, port)
		}
	}
}

// scanCommonPorts is
func scanCommonPorts(host string, wg *sync.WaitGroup) []string {
	defer wg.Done()
	var result []string
	for _, port := range commonPorts {
		ok, err := scanPort(host, port)
		if err != nil {
			continue
		}

		if ok {
			fmt.Printf("[*] FOUND %s:%s\n", host, port)
			result = append(result, fmt.Sprintf("%s:%s", host, port))
		}
	}

	return result
}

func scanPort(host, port string) (bool, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 10*time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// fmt.Printf("[*] connected to %s on port %s\n", host, port)
	// fmt.Printf("[*] sending data to %s on port %s\n", host, port)
	// conn.Write([]byte("hello from world"))
	// buf := make([]byte, 24)
	// conn.Read(buf)
	// fmt.Println("[*] messaged recevied back:")
	// fmt.Println("[*]", string(buf))
	// fmt.Print("\n\n")

	return true, nil
}
