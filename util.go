package botnet

import (
	"crypto/rand"
	"fmt"
)

// Msg is
func Msg(msg ...interface{}) {
	m := append([]interface{}{"[*]"}, msg...)
	fmt.Println(m...)
}

// Err is
func Err(msg ...interface{}) {
	m := append([]interface{}{"[ERROR]"}, msg...)
	fmt.Println(m...)
}

// Debug is
func Debug(msg ...interface{}) {
	m := append([]interface{}{"[DEBUG]"}, msg...)
	fmt.Println(m...)
}

// GenerateKey generates a 32 byte key
func GenerateKey() ([]byte, error) {
	key := make([]byte, 64)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
