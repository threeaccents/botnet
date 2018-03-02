package botnet

//PortScanner is
type PortScanner interface {
	Scan(host string, ports []string) <-chan string
	FullScan() <-chan string
	SimpleScan() <-chan string
}
