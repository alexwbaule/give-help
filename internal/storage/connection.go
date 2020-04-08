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
	Db      *sql.DB
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
			ret.Db = db
		}		
	}		

	return ret, err
}

func (c *Connection) Close() {
	if c != nil {
		if c.Db != nil {
			c.Db.Close()
		}
	}
}

func (c *Connection) Execute(cmd string, args ...interface{}) (int, error) {
	stmt, err := c.Db.Prepare(cmd)

	if err != nil {
		log.Error("fail to try prepare sql command: %s", cmd)
	}

	res, err = stmt.Exec(args)

	if err != nil {
		log.Error("fail to try execute sql command: %s", cmd)
	}

	return res.RowsAffected()
}