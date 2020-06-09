package connection

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

//Connection Object struct
type Connection struct {
	config *Config
	Client *elasticsearch.Client
}

type Config struct {
	Addresses []string
}

//New creates a new instance
func New(cfg *Config) (*Connection, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
	})

	if err != nil {
		log.Printf("error creating the elastic search client: %s", err)
		return nil, err
	} else {
		info, _ := es.Info()
		log.Printf("Elastic cache connected: %s\n", info)
	}

	return &Connection{
		config: cfg,
		Client: es,
	}, err
}
