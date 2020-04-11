package transaction

import (
	"fmt"
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
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
}

func getGiverID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
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

func prepare() (*Transaction, *models.Transaction) {
	trs := createTransaction()
	service := createHandler()

	return service, trs
}

func TestInsert(t *testing.T) {
	service, trs := prepare()

	id, err := service.Insert(trs)

	if err != nil {
		t.Errorf("fail to insert transaction data from %v - error: %s", *trs, err.Error())
	}

	if len(id) == 0 {
		t.Errorf("fail to try insert transaction data from %v - error: %s", *trs, fmt.Errorf("empty user id on return"))
	}
}

func TestLoad(t *testing.T) {
	service, trs := prepare()

	id, err := service.Insert(trs)

	if err != nil {
		t.Errorf("fail to insert transaction data from %v - error: %s", *trs, err.Error())
	}

	if len(id) == 0 {
		t.Errorf("fail to try insert transaction data from %v - error: %s", *trs, fmt.Errorf("empty user id on return"))
	}

	if id != trs.TransactionID {
		t.Errorf("fail to try load transaction. Expected: %s Received: %s", trs.TransactionID, id)
	}
}

func TestInsertReview(t *testing.T) {
	service, trs := prepare()

	id, err := service.Insert(trs)

	if err != nil {
		t.Errorf("fail to insert transaction data from %v - error: %s", *trs, err.Error())
	}

	if len(id) == 0 {
		t.Errorf("fail to try insert transaction data from %v - error: %s", *trs, fmt.Errorf("empty user id on return"))
	}

	giverReview := &models.Review{
		Comment: "Curti o cara, gente fina!",
		Rating:  10,
	}

	err = service.InsertGiverReview(string(trs.TransactionID), giverReview)

	if err != nil {
		t.Errorf("fail to insert transaction review from %v - error: %s", *trs, err.Error())
	}

	trs, err = service.Load(string(trs.TransactionID))

	if err != nil {
		t.Errorf("fail to load transaction data from %v - error: %s", trs.TransactionID, err.Error())
	}

	if trs.GiverReview.Comment != giverReview.Comment {
		t.Errorf("fail to update transaction review (giver comment) from %v - error: %s", trs.TransactionID, err.Error())
	}

	if trs.GiverReview.Rating != giverReview.Rating {
		t.Errorf("fail to update transaction review (giver rating) from %v - error: %s", trs.TransactionID, err.Error())
	}

	takerReview := &models.Review{
		Comment: "Meio diferente, mas foi muito legal!",
		Rating:  9,
	}

	err = service.InsertTakerReview(string(trs.TransactionID), takerReview)

	if err != nil {
		t.Errorf("fail to insert transaction review from %v - error: %s", *trs, err.Error())
	}

	trs, err = service.Load(string(trs.TransactionID))

	if err != nil {
		t.Errorf("fail to load transaction data from %v - error: %s", trs.TransactionID, err.Error())
	}

	if trs.TakerReview.Comment != takerReview.Comment {
		t.Errorf("fail to update transaction review (taker comment) from %v - error: %s", trs.TransactionID, err.Error())
	}

	if trs.TakerReview.Rating != takerReview.Rating {
		t.Errorf("fail to update transaction review (taker rating) from %v - error: %s", trs.TransactionID, err.Error())
	}
}

func TestChangeTransactionStatus(t *testing.T) {
	service, trs := prepare()

	id, err := service.Insert(trs)

	if err != nil {
		t.Errorf("fail to insert transaction data from %v - error: %s", *trs, err.Error())
	}

	if len(id) == 0 {
		t.Errorf("fail to try insert transaction data from %v - error: %s", *trs, fmt.Errorf("empty user id on return"))
	}

	err = service.ChangeTransactionStatus(string(trs.TransactionID), models.TransactionStatusDone)

	if err != nil {
		t.Errorf("fail to update transaction status to %v - error: %s", models.TransactionStatusDone, err.Error())
	}

	trs, err = service.Load(string(trs.TransactionID))

	if err != nil {
		t.Errorf("fail to load transaction data from %v - error: %s", trs.TransactionID, err.Error())
	}

	if trs.Status != models.TransactionStatusDone {
		t.Errorf("invalid status from transaction. Expected: %s Received: %s", models.TransactionStatusDone, trs.Status)
	}
}
