package ssh

import (
	"fmt"
	"os"
	"os/exec"
)

func scp(filepath, destpath string, cred *credential) error {
	fmt.Println("[*] sending file to", cred.host)
	// Open up the file to be sent over
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get the file stats
	s, err := f.Stat()
	if err != nil {
		return err
	}

	// Check if its a file or directory
	if s.IsDir() {
		return scpDir(filepath, destpath, cred)
	}

	return scpFile(filepath, destpath, cred)
}

func scpDir(filepath, destpath string, cred *credential) error {
	cmd := exec.Command("scp", "-r", filepath, cred.username+"@"+cred.host+":"+destpath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func scpFile(filepath, destpath string, cred *credential) error {
	cmd := exec.Command("sshpass -p "+cred.password+" scp", filepath, cred.username+"@"+cred.host+":"+destpath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
