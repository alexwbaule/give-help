package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/prometheus/common/log"
)

type Config struct {
	Host   string
	User   string
	Pass   string
	DBName string
}

type Connection struct {
	config  *Config
	connStr string
}

func New(config *Config) *Connection {
	ret := &Connection{
		config: config,
	}

	ret.connStr = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		config.User,
		config.Pass,
		config.DBName,
		config.Host,
	)

	return ret
}

func (c *Connection) Get() *sql.DB {
	db, err := sql.Open("postgres", c.connStr)

	if err == nil {
		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}
	}

	return db
}
