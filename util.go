package botnet

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

//Bytes transform what ever payload is past in to a byte array
func Bytes(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

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
