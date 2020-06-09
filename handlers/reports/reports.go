package reports

import (
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/proposal"
)

//Reports Object struct
type Reports struct {
	storage *storage.Proposal
}

//New creates a new instance
func New(conn *connection.Connection) *Reports {
	return &Reports{
		storage: storage.New(conn),
	}
}

func (r *Reports) LoadViews() ([]*models.ProposalReport, error) {
	ret, err := r.storage.LoadViews()

	if err != nil {
		log.Printf("fail to load proposal views: %s", err)
	}

	return ret, err
}

func (r *Reports) LoadViewsCSV() (string, error) {
	data, err := r.storage.LoadViewsCSV()

	if err != nil {
		log.Printf("fail to load proposal views: %s", err)
	}

	return data, err
}
