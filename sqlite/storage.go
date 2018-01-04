package sqlite

import (
	"database/sql"
	"encoding/hex"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rodzzlessa24/botnet"
)

//Client i
type Client struct {
	DB *sql.DB
}

//Open is
func Open(source string) (*sql.DB, error) {
	return sql.Open("sqlite3", source)
}

//CreateTables is
func (c *Client) CreateTables() error {
	q := `CREATE TABLE IF NOT EXISTS bots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bot_id VARCHAR(32) NOT NULL UNIQUE,
		host VARCHAR(15) NOT NULL,
		port VARCHAR(5) NOT NULL UNIQUE,
		created_at DATE NOT NULL,
		updated_at DATE NOT NULL
	);
	`
	stmt, err := c.DB.Prepare(q)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

//AddBot is
func (c *Client) AddBot(b *botnet.Bot) (*botnet.Bot, error) {
	q := `
		INSERT INTO bots 
			(bot_id, host, port, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?)
	`

	if _, err := c.exec(q, hex.EncodeToString(b.ID), b.Host, b.Port, time.Now(), time.Now()); err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) exec(q string, args ...interface{}) (sql.Result, error) {
	stmt, err := c.DB.Prepare(q)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(args...)
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
	q := `SELECT bot_id, host, port FROM bots`

	rows, err := c.DB.Query(q)
	if err != nil {
		return nil, err
	}

	bots := []*botnet.Bot{}
	for rows.Next() {
		var id string
		var host string
		var port string
		if err := rows.Scan(&id, &host, &port); err != nil {
			return nil, err
		}
		b := new(botnet.Bot)
		b.ID, _ = hex.DecodeString(id)
		b.Host = host
		b.Port = port
		bots = append(bots, b)
	}

	return bots, nil
}

//GetBot is
func (c *Client) GetBot(id []byte) (*botnet.Bot, error) {
	return nil, nil
}
