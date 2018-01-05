package botnet

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

//CommonIPPorts is
var CommonIPPorts = []string{"22"}

//IPPorts is
func IPPorts() []string {
	var ports []string
	for p := 1; p < 9999; p++ {
		ports = append(ports, strconv.Itoa(p))
	}
	return ports
}

//Scanner is
type Scanner struct {
	Ports []string
	Hosts []string

	Results []string

	wg          *sync.WaitGroup
	errCh       chan error
	jobQueue    chan job
	workerQueue chan chan job
	workers     []worker
}

//NewScanner t is the type of scanner.
func NewScanner(hosts, ports []string, maxQueue, maxWorkers int) *Scanner {

	s := &Scanner{
		Ports:       ports,
		Hosts:       hosts,
		wg:          new(sync.WaitGroup),
		errCh:       make(chan error),
		jobQueue:    make(chan job),
		workerQueue: make(chan chan job, maxWorkers),
	}

	for i := 0; i < maxWorkers; i++ {
		fmt.Println("starting worker", i+1)
		worker := newWorker(s.workerQueue)
		s.workers = append(s.workers, worker)
		worker.Start()
	}

	go func() {
		for {
			select {
			case job := <-s.jobQueue:
				go func() {
					worker := <-s.workerQueue
					worker <- job
				}()
			}
		}
	}()

	return s
}

//Scan is
func (s *Scanner) Scan() <-chan string {
	resCh := make(chan string)
	totalScans := len(s.Hosts) * len(s.Ports)
	s.wg.Add(totalScans)

	go func() {
		for err := range s.errCh {
			log.Println(err)
		}
	}()

	go func() {
		// loop over hosts
		for _, h := range s.Hosts {
			go func(host string) {
				for _, p := range s.Ports {
					job := job{
						resCh:   resCh,
						addr:    fmt.Sprintf("%s:%s", host, p),
						scanner: s,
					}
					s.jobQueue <- job
				}
			}(h)
		}
		s.wg.Wait()
		s.stopWorkers()
		close(resCh)
	}()

	return resCh
}

func (s *Scanner) stopWorkers() {
	for _, w := range s.workers {
		w.Stop()
	}
}

func (s *Scanner) scanPort(addr string) error {
	defer s.wg.Done()
	conn, err := net.DialTimeout("tcp", addr, 15*time.Second)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return ErrTimeout
		}
		return err
	}
	conn.Close()
	return nil
}

//job is
type job struct {
	resCh   chan string
	addr    string
	scanner *Scanner
}

type worker struct {
	Work        chan job
	WorkerQueue chan chan job
	QuitChan    chan bool
}

func newWorker(workerQueue chan chan job) worker {
	return worker{
		Work:        make(chan job),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
}

func (w worker) Start() {
	go func() {
		for {
			// alert the workque we have a worker available
			w.WorkerQueue <- w.Work

			select {
			case job := <-w.Work:
				time.Sleep(1 * time.Second)
				if err := job.scanner.scanPort(job.addr); err != nil {
					if err != ErrTimeout {
						job.scanner.errCh <- err
					}
					break
				}
				job.resCh <- job.addr

			case <-w.QuitChan:
				// We have been asked to stop.
				close(w.Work)
				fmt.Printf("worker stopping\n")
				return
			}
		}
	}()
}

func (w worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
