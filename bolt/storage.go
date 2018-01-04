package bolt

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/rodzzlessa24/botnet"
)

//Client i
type Client struct {
	DB *bolt.DB
}

//Open is
func Open(path string, mode os.FileMode, options *bolt.Options) (*bolt.DB, error) {
	return bolt.Open(path, mode, options)
}

//AddBot is
func (c *Client) AddBot(b *botnet.Bot) (*botnet.Bot, error) {
	return nil, nil
}

//RemoveBot is
func (c *Client) RemoveBot(id []byte) error {
	return nil
}

//ListBots is
func (c *Client) ListBots() ([]*botnet.Bot, error) {
	return nil, nil
}

//GetBot is
func (c *Client) GetBot(id []byte) (*botnet.Bot, error) {
	return nil, nil
}
