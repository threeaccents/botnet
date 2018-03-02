package botnet

//PortScanner is
type PortScanner interface {
	Scan() <-chan string
}
