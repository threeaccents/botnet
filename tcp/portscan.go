package tcp

import "github.com/rodzzlessa24/botnet/lib/portscan"

//PortScanService is
type PortScanService struct {
}

//Scan is
func (p *PortScanService) Scan(hosts, ports []string) <-chan string {
	scanner := portscan.New(hosts, ports, 50)
	return scanner.Scan()
}

//FullScan is
func (p *PortScanService) FullScan(hosts []string) <-chan string {
	scanner := portscan.New(hosts, portscan.IPPorts(), 50)
	return scanner.Scan()
}

//SimpleScan is
func (p *PortScanService) SimpleScan(hosts []string) <-chan string {
	scanner := portscan.New(hosts, portscan.CommonIPPorts, 50)
	return scanner.Scan()
}
