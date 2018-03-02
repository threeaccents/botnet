package botnet

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

//Botter is
type Botter interface {
	RansomCompleteCmd()

	HandleScan()
	HandleRansomware()
}

//Bot is
type Bot struct {
	ID     []byte
	Host   string
	Port   string
	CCAddr string
}

//Bytes is
func (b *Bot) Bytes() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := gob.NewEncoder(buff).Encode(b); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

//Addr is
func (b *Bot) Addr() string {
	return fmt.Sprintf("%s:%s", b.Host, b.Port)
}

//BytesToBot is
func BytesToBot(b []byte) (*Bot, error) {
	bot := new(Bot)
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(bot); err != nil {
		return nil, err
	}
	return bot, nil
}
