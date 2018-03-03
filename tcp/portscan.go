package tcp

//PortScanService is
type PortScanService struct {
}

//Scan is
func (p *PortScanService) Scan(hosts, ports []string) <-chan string {
	return make(chan string)
}

//FullScan is
func (p *PortScanService) FullScan(hosts []string) <-chan string {
	return make(chan string)
}

//SimpleScan is
func (p *PortScanService) SimpleScan(hosts []string) <-chan string {
	return make(chan string)
}
