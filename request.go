package botnet

import (
	"bytes"
	"encoding/gob"
)

//RansomCompleteRequest is
type RansomCompleteRequest struct {
	BotID []byte
	Key   []byte
}

//Bytes is
func Bytes(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
