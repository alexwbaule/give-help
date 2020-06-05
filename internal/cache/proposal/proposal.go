package proposal

import (
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/elastic/go-elasticsearch/v8"
)

//Proposal Object struct
type Proposal struct {
	config elasticsearch.Config
	client *elasticsearch.Client
}

//New creates a new instance
func New(cfg elasticsearch.Config) (*Proposal, error) {
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Printf("error creating the elastic search client: %s", err)
		return nil, err
	} else {
		log.Println(es.Info())
	}

	return &Proposal{
		config: cfg,
		client: es,
	}, err
}

//Upsert insert or update
func (p *Proposal) Upsert(proposal *models.Proposal) error {
	var err error
	return err
}

//Find find all proposals that match with filter
func (p *Proposal) Find(filter *models.Filter) ([]*models.Proposal, error) {
	return []*models.Proposal{}, nil
}
