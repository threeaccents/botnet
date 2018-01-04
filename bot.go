package botnet

import (
	"bytes"
	"encoding/gob"
)

//Bot is
type Bot struct {
	ID   []byte
	Host string
	Port string
}

//Bytes is
func (b *Bot) Bytes() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := gob.NewEncoder(buff).Encode(b); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

//BytesToBot is
func BytesToBot(b []byte) (*Bot, error) {
	bot := new(Bot)
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(bot); err != nil {
		return nil, err
	}
	return bot, nil
}
