package connection

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/prometheus/common/log"
)

//Connection base connection struct
type Connection struct {
	config   *common.DbConfig
	connStr  string
	db       *sql.DB
	lastPing time.Time
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

	ret.createConnection()

	return ret
}

//Get returns an open sql.DB object (this must be closed after use!)
func (c *Connection) Get() *sql.DB {
	if c.db == nil {
		c.createConnection()
	}

	if time.Since(c.lastPing) > time.Minute {
		log.Print("sending ping to database")
		if err := c.db.Ping(); err != nil {
			log.Print("ping to database fail, trying to reconnect")
			c.createConnection()
		}
		c.lastPing = time.Now()
	}

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

func (c *Connection) createConnection() {
	db, err := sql.Open("postgres", c.connStr)

	if err == nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	c.lastPing = time.Now()

	log.Print("db connection created")
	c.db = db
}
