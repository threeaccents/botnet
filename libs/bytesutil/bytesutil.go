package bytesutil

import (
	"bytes"
	"encoding/gob"
)

// Unmarshal will transform a byte array into the passed in value.
func Unmarshal(b []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(b)).Decode(v)
}

// Marshal will transform a byte array into the passed in value.
func Marshal(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
