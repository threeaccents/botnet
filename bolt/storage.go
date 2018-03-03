package bolt

import (
	bolt "github.com/coreos/bbolt"
	"github.com/rodzzlessa24/botnet"
)

const (
	botsBucket = "bots"
)

//Client is
type Client struct {
	DB *bolt.DB
}

// Open creates and opens a database at the given path.
// If the file does not exist then it will be created automatically.
// Passing in nil options will cause Bolt to open the database with the default options.
func Open(source string) (*bolt.DB, error) {
	return bolt.Open(source, 0600, nil)
}

//NewClient is
func NewClient(db *bolt.DB) (*Client, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(botsBucket))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &Client{db}, nil
}

//AddBot is
func (c *Client) AddBot(b *botnet.Bot) (*botnet.Bot, error) {
	bytesBot, err := b.Bytes()
	if err != nil {
		return nil, err
	}
	if err := c.DB.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte(botsBucket))
		return bu.Put(b.ID, bytesBot)
	}); err != nil {
		return nil, err
	}
	return b, nil
}

//UpdateBot is
func (c *Client) UpdateBot(b *botnet.Bot) (*botnet.Bot, error) {
	bytesBot, err := b.Bytes()
	if err != nil {
		return nil, err
	}
	if err := c.DB.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte(botsBucket))
		return bu.Put(b.ID, bytesBot)
	}); err != nil {
		return nil, err
	}
	return b, nil
}

//AddRansomKey is
func (c *Client) AddRansomKey(botID, key []byte) error {
	return nil
}

//RemoveBot is
func (c *Client) RemoveBot(id []byte) error {
	return nil
}

//ListBots is
func (c *Client) ListBots() ([]*botnet.Bot, error) {
	var bots []*botnet.Bot
	if err := c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(botsBucket))
		return b.ForEach(func(_, v []byte) error {
			botStruct, err := botnet.BytesToBot(v)
			if err != nil {
				return err
			}
			bots = append(bots, botStruct)
			return nil
		})
	}); err != nil {
		return nil, err
	}

	return bots, nil
}

//GetBot is
func (c *Client) GetBot(id []byte) (*botnet.Bot, error) {
	return nil, nil
}
