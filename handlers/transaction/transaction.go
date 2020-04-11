package transaction

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	conn "github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/transaction"
)

//User Object struct
type transaction struct {
	storage *storage.transaction
	config  *common.Config
}

//New creates a new instance
func New(config *common.Config) *transaction {
	conn := conn.New(config.Db)
	return &transaction{
		storage: storage.New(conn),
		config:  config,
	}
}

func (t *transaction) Insert(transaction *models.Transaction) (models.ID, error) {

}

func (t *transaction) update(transaction *models.Transaction) error {
}

func (t *transaction) LoadByProposalID(proposalID string) ([]*models.Transaction, error) {

}

func (t *transaction) LoadByUserID(userID string) ([]*models.Transaction, error) {

}

func (t *transaction) InsertGiverReview(transactionID string, review *models.Review) error {

}

func (t *transaction) InsertTakerReview(transactionID string, review *models.Review) error {

}

func (t *transaction) ChangeTransactionStatus(transactionID string, newStatus *models.TransactionStatus) error {

}
