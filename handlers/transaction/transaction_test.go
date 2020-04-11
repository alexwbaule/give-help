package transaction

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
)

func createHandler() *Transaction {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	config := &common.Config{
		Db: dbConfig,
	}

	return New(config)
}

func getTakerID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GGG"
}

func getGiverID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8TIT"
}

func getUserID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8UUU"
}

func getProposalID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8PPP"
}

func getTransactionID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8TTT"
}

func createTransaction() *models.Transaction {
	return &models.Transaction{
		TransactionID: models.ID(getTransactionID()),
		ProposalID:    models.ID(getProposalID()),
		GiverID:       models.ID(getGiverID()),
		TakerID:       models.ID(getTakerID()),
		GiverReview: &models.Review{
			Rating:  0,
			Comment: "",
		},
		TakerReview: &models.Review{
			Rating:  0,
			Comment: "",
		},
		Status: models.TransactionStatusOpen,
	}
}

func TestInsert(t *testing.T) {

}

func TestUpdate(t *testing.T) {

}

func TestLoad(t *testing.T) {

}

func TestInsertReview(t *testing.T) {

}

func TestChangeTransactionStatus(t *testing.T) {

}
