package botnet

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

//CommonIPPorts is
var CommonIPPorts = []string{"22"}

//IPPorts is
func IPPorts() []string {
	var ports []string
	for p := 1; p < 999; p++ {
		ports = append(ports, strconv.Itoa(p))
	}
	return ports
}

//Scanner is
type Scanner struct {
	Ports []string
	Hosts []string

	wg    *sync.WaitGroup
	errCh chan error
	sem   chan int
}

//NewScanner t is the type of scanner.
func NewScanner(hosts, ports []string, maxQueue int) *Scanner {
	s := &Scanner{
		Ports: ports,
		Hosts: hosts,
		wg:    new(sync.WaitGroup),
		errCh: make(chan error),
		sem:   make(chan int, maxQueue),
	}

	return s
}

//Scan is
func (s *Scanner) Scan() <-chan string {
	resCh := make(chan string)
	totalScans := len(s.Hosts) * len(s.Ports)
	s.wg.Add(totalScans)

	// listen for errors
	go func() {
		for err := range s.errCh {
			if e, ok := err.(net.Error); ok && !e.Timeout() && !strings.Contains(err.Error(), "getsockopt: connection refused") {
				log.Println(err)
			}
		}
	}()

	go func() {
		// loop over hosts
		for _, h := range s.Hosts {
			go func(host string) {
				for _, p := range s.Ports {
					addr := fmt.Sprintf("%s:%s", host, p)
					s.sem <- 1
					go s.scanPort(addr, resCh)
				}
			}(h)
		}
		s.wg.Wait()
		close(s.errCh)
		close(resCh)
	}()

	return resCh
}

func (s *Scanner) scanPort(addr string, resCh chan string) {
	defer s.wg.Done()
	conn, err := net.DialTimeout("tcp", addr, 15*time.Second)
	if err != nil {
		s.errCh <- err
		<-s.sem
		return
	}
	conn.Close()
	resCh <- addr
	<-s.sem
}
