package storage

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/prometheus/common/log"
)

//Config base connection config struct
type Config struct {
	Host   string
	User   string
	Pass   string
	DBName string
}

//Connection base connection struct
type Connection struct {
	config  *Config
	connStr string
}

//New creates a connection helper
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

//Get returns an open sql.DB object (this must be closed after use!)
func (c *Connection) Get() *sql.DB {
	db, err := sql.Open("postgres", c.connStr)

	if err == nil {
		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}
	}

	return db
}

//GetMessage try get better error message
func (c *Connection) CheckError(err error) error {
	if perr, ok := err.(*pq.Error); ok {
		return fmt.Errorf("%s - pq-error: %v", perr.Error(), perr)
	}

	return err
}
