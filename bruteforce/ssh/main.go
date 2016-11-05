package ssh

import "golang.org/x/crypto/ssh"

// Attack check if the give username and password combo are valid credentials
func Attack(username, password string) (bool, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	_, err := ssh.Dial("tcp", "192.168.0.6:22", config)
	if err != nil {
		return false, err
	}

	return true, nil
}
