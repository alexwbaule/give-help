package proposal

import (
	"bytes"
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

	//log.Printf("updating cache for: %s", string(raw))

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
			log.Printf("error parsing the response body: %s", err)
		} else {
			log.Printf("cache proposal index is updated [id=%s]: [%s] %s; version=%d", proposal.ProposalID, res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return err
}

//Find find all proposals that match with filter
func (p *Proposal) Find(filter *models.Filter) ([]*models.Proposal, error) {
	args := []interface{}{}

	if len(filter.Description) > 0 {
		args = append(args, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"fields": []string{
					"title^3",
					"description^2",
					"tags^2",
					"target_area.area_tags",
					"target_area.city",
					"target_area.country",
					"target_area.state",
				},
				"type":  "best_fields",
				"query": filter.Description,
			},
		})
	}

	if len(filter.UserID) > 0 {
		args = append(args, map[string]interface{}{
			"match": map[string]interface{}{
				"user_id": filter.UserID,
			},
		})
	}

	if len(filter.Side) > 0 {
		args = append(args, map[string]interface{}{
			"match": map[string]interface{}{
				"side": filter.Side,
			},
		})
	}

	for _, t := range filter.ProposalTypes {
		args = append(args, map[string]interface{}{
			"match": map[string]interface{}{
				"proposal_type": t,
			},
		})
	}

	for _, t := range filter.Tags {
		args = append(args, map[string]interface{}{
			"match": map[string]interface{}{
				"tags": t,
			},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": args,
			},
		},
	}

	/* to debug ES query
	raw, _ := json.MarshalIndent(query, "", "\t")
	log.Printf("QUERY => %s", string(raw))
	*/
	return p.load(query)
}

func (p *Proposal) LoadFromID(proposalID string) (*models.Proposal, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"proposal_id": proposalID,
			},
		},
	}

	ret, err := p.load(query)

	if err != nil {
		log.Printf("fail to find this proposal on cache: %s", err)
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("proposal not found on cache")
	}

	return ret[0], err
}

func (p *Proposal) load(query map[string]interface{}) ([]*models.Proposal, error) {
	var ret []*models.Proposal
	var err error

	client := p.conn.Client

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return ret, err
	}

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(proposalIndexName),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)

	if err != nil {
		log.Printf("fail to try search documents on cache: %s", err)
		return ret, err
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			err = fmt.Errorf("error searching documents [%s]", res.Status())
			log.Printf("fail to try search documents on cache: %s", err)
		} else {
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return ret, err
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("error parsing the response body: %s", err)
		return ret, err
	}

	raw, err := json.Marshal(r["hits"].(map[string]interface{})["hits"])

	if err != nil {
		log.Printf("fail to try read result from cache: %s", err)
		return ret, err
	}

	var results []esResult

	if err = json.Unmarshal(raw, &results); err != nil {
		log.Printf("error to try parse source proposal from cache: %s", err)
	} else {
		for _, p := range results {
			ret = append(ret, &p.Source)
		}
	}

	return ret, err
}

type esResult struct {
	Index  string          `json:"_index,omitempty"`
	Score  float64         `json:"_score,omitempty"`
	Source models.Proposal `json:"_source,omitempty"`
}

//Reindex refresh index
func (p *Proposal) Reindex(proposals []*models.Proposal) {
	var err error
	for _, item := range proposals {
		if item != nil {
			err = p.Upsert(item)

			if err != nil {
				log.Printf("fail to try reindex proposal [id:%s]: %s\n", item.ProposalID, err)
			}
		}
	}
}
