package proposal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/cache/connection"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const proposalIndexName = "proposals"

//Proposal Object struct
type Proposal struct {
	conn *connection.Connection
}

//New creates a new instance
func New(conn *connection.Connection) *Proposal {
	return &Proposal{conn: conn}
}

//Upsert insert or update index data
func (p *Proposal) Upsert(proposal *models.Proposal) error {
	var err error

	raw, err := json.MarshalIndent(proposal, "", "\t")

	req := esapi.IndexRequest{
		Index:      proposalIndexName,
		DocumentID: string(proposal.ProposalID),
		Body:       strings.NewReader(string(raw)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), p.conn.Client)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		err = fmt.Errorf("error indexing document ID=%s [%s]", proposal.ProposalID, res.Status())
		log.Printf("fail to try insert proposal on cache: %s", err)
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			log.Printf("ES proposal index is updated: [%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return err
}

//Find find all proposals that match with filter
func (p *Proposal) Find(filter *models.Filter) ([]*models.Proposal, error) {
	return []*models.Proposal{}, nil
}
