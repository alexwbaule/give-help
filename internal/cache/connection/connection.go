package connection

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

//Connection Object struct
type Connection struct {
	config *elasticsearch.Config
	Client *elasticsearch.Client
}

//New creates a new instance
func New(conn *elasticsearch.Config) (*Connection, error) {
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Printf("error creating the elastic search client: %s", err)
		return nil, err
	} else {
		log.Println(es.Info())
	}

	return &Connection{
		config: cfg,
		Client: es,
	}, err
}
