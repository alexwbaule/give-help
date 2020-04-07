package storage

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"database/sql"

	_ "github.com/lib/pq"
)

type Config struct {
	Host string
	User string
	Pass string
	DBName string
}

type Connection struct {
	config *Config
	db      *sql.DB
}

func New(config *Config) (*Connection, error) {
	ret := &Storage{
		config: config
	}
	
	connStr := ftm.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=verify-full", 
		config.User,
		config.Pass,
		config.DBName,
		config.Host,
	)

	db, err := sql.Open("postgres", connStr)
	
	if err == nil {
		err = db.Ping()
		if err = nil {
			ret.db = db
		}		
	}		

	return ret, err
}

func (c *Connection) Close() {
	if c != nil {
		if c.DB != nil {
			c.DB.Close()
		}
	}
}