package botnet

//PortScanner is
type PortScanner interface {
	Scan(hosts, ports []string) <-chan string
	FullScan(hosts []string) <-chan string
	SimpleScan(host []string) <-chan string
}
