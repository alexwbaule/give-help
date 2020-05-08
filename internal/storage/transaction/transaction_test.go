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

func insert(t *testing.T) {
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

	loaded, err = storage.LoadByUserID(string(loaded[0].GiverID))

	if err != nil {
		t.Errorf("fail to try load transaction data from giverID: %s - error: %s", data.ProposalID, err)
	}

	loaded, err = storage.LoadByUserID(string(loaded[0].TakerID))

	if err != nil {
		t.Errorf("fail to try load transaction data from takerID: %s - error: %s", data.ProposalID, err)
	}
}

func update(t *testing.T) {
	storage := createConn()
	initial := createTransaction()

	err := storage.Upsert(initial)

	data, err := storage.LoadByProposalID(string(initial.ProposalID))

	if err != nil {
		t.Errorf("fail to try load transaction data from proposalID: %s - error: %s", initial.ProposalID, err)
	}

	trs := data[0]

	trs.Status = models.TransactionStatusInProgress
	trs.GiverReview.Rating = 4
	trs.GiverReview.Comment = "Cara gente fina, deu tudo certo!"

	err = storage.Upsert(trs)

	if err != nil {
		t.Errorf("fail to try update transaction data from giverID: %s - error: %s", trs.ProposalID, err)
	}

	data, err = storage.LoadByProposalID(string(trs.ProposalID))

	updated := data[0]

	if err != nil {
		t.Errorf("fail to try load transaction data from giverID: %s - error: %s", trs.ProposalID, err)
	}

	if updated.Status != trs.Status {
		t.Errorf("fail to try update transaction status to inProgress - error: %s expected: %s received: %s", err, trs.Status, updated.Status)
	}

	if updated.GiverReview.Rating != trs.GiverReview.Rating {
		t.Errorf("fail to try update transaction giver review rating - error: %s expected: %s received: %s", err, trs.Status, updated.Status)
	}

	if updated.GiverReview.Comment != trs.GiverReview.Comment {
		t.Errorf("fail to try update transaction giver review comment - error: %s expected: %s received: %s", err, trs.Status, updated.Status)
	}

	trs.Status = models.TransactionStatusDone
	trs.TakerReview.Rating = 5
	trs.TakerReview.Comment = "Foi muito legal, recomendo!"

	err = storage.Upsert(trs)

	if err != nil {
		t.Errorf("fail to try update transaction data from giverID: %s - error: %s", trs.ProposalID, err)
	}

	data, err = storage.LoadByProposalID(string(trs.ProposalID))

	updated = data[0]

	if err != nil {
		t.Errorf("fail to try load transaction data from giverID: %s - error: %s", trs.ProposalID, err)
	}

	if updated.Status != trs.Status {
		t.Errorf("fail to try update transaction status to done - error: %s expected: %s received: %s", err, trs.Status, updated.Status)
	}

	if updated.TakerReview.Rating != trs.TakerReview.Rating {
		t.Errorf("fail to try update transaction taker review rating - error: %s expected: %s received: %s", err, trs.Status, updated.Status)
	}

	if updated.TakerReview.Comment != trs.TakerReview.Comment {
		t.Errorf("fail to try update transaction taker review comment - error: %s expected: %s received: %s", err, trs.Status, updated.Status)
	}
}

func Test(t *testing.T) {
	insert(t)
	update(t)
}
