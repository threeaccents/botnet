package queue

import "fmt"

// We have

//Channel is
type Channel struct {
	jobQueue    chan []byte
	workerQueue chan chan []byte
}

//NewChannel will
func NewChannel(maxQueue, maxWorkers int) *Channel {
	c := &Channel{
		jobQueue:    make(chan []byte, maxQueue),
		workerQueue: make(chan chan []byte, maxWorkers),
	}

	go func() {
		for {
			select {
			case job := <-c.jobQueue:
				fmt.Println(string(job))
			}
		}
	}()

	return c
}

//Push is
func (c *Channel) Push(msg []byte) {
	c.jobQueue <- msg
}

//Consume is
func (c *Channel) Consume() (chan []byte, chan bool) {
	return make(chan []byte), make(chan bool)
}
