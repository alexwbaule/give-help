package connection

import (
	"log"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/elastic/go-elasticsearch/v8"
)

//Connection Object struct
type Connection struct {
	config *common.CacheConfig
	Client *elasticsearch.Client
}

//New creates a new instance
func New(cfg *common.CacheConfig) (*Connection, error) {
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
