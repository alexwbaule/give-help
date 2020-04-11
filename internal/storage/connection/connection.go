package connection

import (
	"database/sql"
	"fmt"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/prometheus/common/log"
)

//Connection base connection struct
type Connection struct {
	config  *common.DbConfig
	connStr string
	db      *sql.DB
}

//New creates a connection helper
func New(config *common.DbConfig) *Connection {
	ret := &Connection{
		config: config,
	}

	ret.connStr = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		config.User,
		config.Pass,
		config.DBName,
		config.Host,
	)

	db, err := sql.Open("postgres", ret.connStr)

	if err == nil {
		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}
	}

	ret.db = db

	return ret
}

//Get returns an open sql.DB object (this must be closed after use!)
func (c *Connection) Get() *sql.DB {
	return c.db
}

func (c *Connection) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

//GetMessage try get better error message
func (c *Connection) CheckError(err error) error {
	if perr, ok := err.(*pq.Error); ok {
		return fmt.Errorf("%s - pq-error: %v", perr.Error(), perr)
	}

	return err
}
