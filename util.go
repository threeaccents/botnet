package botnet

import (
	"fmt"
	"os"
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

// Msg is
func Msg(msg ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		m := append([]interface{}{"[*]"}, msg...)
		fmt.Println(m...)
	}
}

// Err is
func Err(msg ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		m := append([]interface{}{"[ERROR]"}, msg...)
		fmt.Println(m...)
	}
}
