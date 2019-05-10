package ransom

import (
	"github.com/threeaccents/botnet"
)

type RansomService struct {
	CryptoService botnet.CryptoService
	Key           []byte
}

func (s *RansomService) Run() error {
	return nil
}
