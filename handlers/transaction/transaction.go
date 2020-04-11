package transaction

import (
	"fmt"
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	conn "github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/transaction"
)

//User Object struct
type Transaction struct {
	storage *storage.Transaction
	config  *common.Config
}

//New creates a new instance
func New(config *common.Config) *Transaction {
	conn := conn.New(config.Db)
	return &Transaction{
		storage: storage.New(conn),
		config:  config,
	}
}

//Insert insert new data
func (t *Transaction) Insert(transaction *models.Transaction) (models.ID, error) {
	if len(transaction.TransactionID) == 0 {
		transaction.TransactionID = models.ID(common.GetULID())
	}

	err := t.storage.Upsert(transaction)

	if err != nil {
		log.Printf("fail to insert new transaction [%s]: %s", transaction.TransactionID, err)
	}

	return transaction.TransactionID, err
}

//Load load data
func (t *Transaction) Load(transactionID string) (*models.Transaction, error) {
	if len(transactionID) == 0 {
		return nil, fmt.Errorf("transactionID is empty")
	}

	ret, err := t.storage.Load(transactionID)

	if err != nil {
		log.Printf("fail to try load transactions from ID [%s]: %s", transactionID, err)
	}

	return ret, err
}

//LoadByProposalID load data
func (t *Transaction) LoadByProposalID(proposalID string) ([]*models.Transaction, error) {
	if len(proposalID) == 0 {
		return nil, fmt.Errorf("proposalID is empty")
	}

	ret, err := t.storage.LoadByProposalID(proposalID)

	if err != nil {
		log.Printf("fail to try load transactions from proposalID [%s]: %s", proposalID, err)
	}

	return ret, err
}

//LoadByUserID load data
func (t *Transaction) LoadByUserID(userID string) ([]*models.Transaction, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("userId is empty")
	}

	ret, err := t.storage.LoadByUserID(userID)

	if err != nil {
		log.Printf("fail to try load transactions from userID [%s]: %s", userID, err)
	}

	return ret, err
}

//InsertGiverReview insert reviews on transaction
func (t *Transaction) InsertGiverReview(transactionID string, review *models.Review) error {
	if review == nil {
		return fmt.Errorf("review is null")
	}

	trs, err := t.Load(transactionID)

	if err != nil {
		log.Printf("fail to try insert review on transactions ID [%s]: %s", transactionID, err)
		return err
	}

	trs.GiverReview = review

	return t.update(trs)
}

//InsertTakerReview insert reviews on transaction
func (t *Transaction) InsertTakerReview(transactionID string, review *models.Review) error {
	if review == nil {
		return fmt.Errorf("review is null")
	}

	trs, err := t.Load(transactionID)

	if err != nil {
		log.Printf("fail to try insert review on transactions ID [%s]: %s", transactionID, err)
		return err
	}

	trs.TakerReview = review

	return t.update(trs)
}

//ChangeTransactionStatus change transaction status
func (t *Transaction) ChangeTransactionStatus(transactionID string, newStatus *models.TransactionStatus) error {
	if newStatus == nil {
		return fmt.Errorf("newStatus is null")
	}

	trs, err := t.Load(transactionID)

	if err != nil {
		log.Printf("fail to try update transaction status [%s] ID [%s]: %s", *newStatus, transactionID, err)
		return err
	}

	trs.Status = *newStatus

	return t.update(trs)
}

func (t *Transaction) update(transaction *models.Transaction) error {
	if transaction == nil {
		return fmt.Errorf("transaction is null")
	}

	err := t.storage.Upsert(transaction)

	if err != nil {
		log.Printf("fail to try insert review on transactions ID [%s]: %s", transaction.TransactionID, err)
	}

	return err
}
