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
	ID   []byte
	Host string
	Port string
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
