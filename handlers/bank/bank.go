package bank

import (
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/bank"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

//Banks Object struct
type Bank struct {
	storage *storage.Bank
}

//New creates a new instance
func New(conn *connection.Connection) *Bank {
	return &Bank{
		storage: storage.New(conn),
	}
}

//Load load banks
func (b *Bank) LoadBanks() ([]*models.Bank, error) {
	ret, err := b.storage.LoadBanks()

	if err != nil {
		log.Printf("fail to try load banks: %s", err)
	}

	return ret, err
}
