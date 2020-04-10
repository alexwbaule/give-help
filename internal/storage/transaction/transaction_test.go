package transaction

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func createConn() *Transaction {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := connection.New(dbConfig)

	return New(conn)
}

func getTakerID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
}

func getGiverID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
}

func getProposalID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GXX"
}

func getTransactionID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GZZ"
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
	storage := createConn()
	data := createTransaction()

	err := storage.Upsert(data)

	if err != nil {
		t.Errorf("fail to try insert transaction data on database - error: %s", err)
	}

	loaded, err := storage.LoadByProposalID(string(data.ProposalID))

	if err != nil {
		t.Errorf("fail to try load transaction data from proposalID: %s - error: %s", data.ProposalID, err)
	}

	loaded, err = storage.LoadByUserID(string(loaded.GiverID))

	if err != nil {
		t.Errorf("fail to try load transaction data from giverID: %s - error: %s", data.ProposalID, err)
	}

	loaded, err = storage.LoadByUserID(string(loaded.TakerID))

	if err != nil {
		t.Errorf("fail to try load transaction data from takerID: %s - error: %s", data.ProposalID, err)
	}
}

func TestUpdate(t *testing.T) {
	storage := createConn()
	data := createTransaction()

	data, err := storage.LoadByProposalID(string(getProposalID()))

	if err != nil {
		t.Errorf("fail to try load transaction data from proposalID: %s - error: %s", data.ProposalID, err)
	}

	data.Status = models.TransactionStatusInProgress
	data.GiverReview.Rating = 4
	data.GiverReview.Comment = "Cara gente fina, deu tudo certo!"

	err = storage.Upsert(data)

	if err != nil {
		t.Errorf("fail to try update transaction data from giverID: %s - error: %s", data.ProposalID, err)
	}

	updated, err := storage.LoadByProposalID(string(data.ProposalID))

	if err != nil {
		t.Errorf("fail to try load transaction data from giverID: %s - error: %s", data.ProposalID, err)
	}

	if updated.Status != data.Status {
		t.Errorf("fail to try update transaction status to inProgress - error: %s expected: %s received: %s", err, data.Status, updated.Status)
	}

	if data.GiverReview.Rating != data.GiverReview.Rating {
		t.Errorf("fail to try update transaction giver review rating - error: %s expected: %s received: %s", err, data.Status, updated.Status)
	}

	if data.GiverReview.Comment != data.GiverReview.Comment {
		t.Errorf("fail to try update transaction giver review comment - error: %s expected: %s received: %s", err, data.Status, updated.Status)
	}

	data.Status = models.TransactionStatusDone
	data.TakerReview.Rating = 5
	data.TakerReview.Comment = "Foi muito legal, recomendo!"

	err = storage.Upsert(data)

	if err != nil {
		t.Errorf("fail to try update transaction data from giverID: %s - error: %s", data.ProposalID, err)
	}

	updated, err = storage.LoadByProposalID(string(data.ProposalID))

	if err != nil {
		t.Errorf("fail to try load transaction data from giverID: %s - error: %s", data.ProposalID, err)
	}

	if updated.Status != data.Status {
		t.Errorf("fail to try update transaction status to done - error: %s expected: %s received: %s", err, data.Status, updated.Status)
	}

	if data.TakerReview.Rating != data.TakerReview.Rating {
		t.Errorf("fail to try update transaction taker review rating - error: %s expected: %s received: %s", err, data.Status, updated.Status)
	}

	if data.TakerReview.Comment != data.TakerReview.Comment {
		t.Errorf("fail to try update transaction taker review comment - error: %s expected: %s received: %s", err, data.Status, updated.Status)
	}
}
