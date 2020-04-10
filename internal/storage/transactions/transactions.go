package transactions

import (
	"fmt"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage"
	"github.com/lib/pq"
)

//Transactions Object struct
type Transactions struct {
	conn *storage.Connection
}

//New creates a new instance
func New(conn *storage.Connection) *Transactions {
	return &Transactions{conn: conn}
}

const upsertTransaction string = `
INSERT INTO TRANSACTIONS
(
	TransactionID,
	ProposalID,
	GiverID,
	TakerID,

	GiverRating,
	GiverReviewComment,
	TakerRating,
	TakerReviewComment,
	Status
)
VALUES
(
	$1, --TransactionID,
	$2, --ProposalID,
	$3, --GiverID,
	$4, --TakerID,

	$5, --GiverRating,
	$6, --GiverReviewComment,
	$7, --TakerRating,
	$8, --TakerReviewComment,
	
	$9 --Status
)
ON CONFLICT (TransactionID) 
DO
	UPDATE
	SET 
		LastUpdate = CURRENT_TIMESTAMP,
		GiverRating = $5,
		GiverReviewComment = $6,
		TakerRating = $7,
		TakerReviewComment = $8,
		
		Status = $9
;`

//Upsert insert or update on database
func (t *Transactions) Upsert(transaction *models.Transaction) error {
	if transaction == nil {
		return fmt.Errorf("cannot insert an empty transaction struct")
	}

	if len(transaction.TransactionID) == 0 {
		return fmt.Errorf("cannot insert an empty TransactionID")
	}

	if len(transaction.ProposalID) == 0 {
		return fmt.Errorf("cannot insert an empty ProposalID")
	}

	if len(transaction.GiverID) == 0 {
		return fmt.Errorf("cannot insert an empty GiverID")
	}

	if len(transaction.TakerID) == 0 {
		return fmt.Errorf("cannot insert an empty TakerID")
	}

	if transaction.GiverReview == nil {
		transaction.GiverReview = &models.Review{}
	}

	if transaction.TakerReview == nil {
		transaction.GiverReview = &models.Review{}
	}

	db := t.conn.Get()
	defer db.Close()

	_, err := db.Exec(
		upsertTransaction,
		transaction.TransactionID,
		transaction.ProposalID,
		transaction.GiverID,
		transaction.TakerID,
		transaction.GiverReview.Rating,
		transaction.GiverReview.Comment,
		transaction.TakerReview.Rating,
		transaction.TakerReview.Comment,
		transaction.Status,
	)

	if err != nil {
		if perr, ok := err.(*pq.Error); ok {
			return fmt.Errorf("fail to try execute upsert transaction data: transaction=%v pq-error=%s", transaction, perr)
		}

		return fmt.Errorf("fail to try execute upsert transaction data: transaction=%v error=%s", transaction, err)
	}

	return nil
}

const selectTransaction string = `
SELECT
	TransactionID,
	ProposalID,
	GiverID,
	TakerID,
	CreatedAt,
	LastUpdate,

	GiverRating,
	GiverReviewComment,
	TakerRating,
	TakerReviewComment,
	Status
FROM
	TRANSACTIONS
WHERE
	%s
ORDER BY
	CreatedAt ASC
`

func (t *Transactions) LoadByProposalID(proposalID string) (*models.Transaction, error) {
	return t.load(fmt.Sprintf(selectTransaction, "ProposalID = $1"), proposalID)
}

func (t *Transactions) LoadByUserID(userID string) (*models.Transaction, error) {
	return t.load(fmt.Sprintf(selectTransaction, "GiverID = $1 OR TakerID = $1"), userID)
}

func (t *Transactions) load(cmd string, args ...interface{}) (*models.Transaction, error) {
	ret := models.Transaction{
		GiverReview: &models.Review{},
		TakerReview: &models.Review{},
	}

	db := t.conn.Get()
	defer db.Close()

	row := db.QueryRow(cmd, args...)

	err := row.Scan(
		&ret.TransactionID,
		&ret.ProposalID,
		&ret.GiverID,
		&ret.TakerID,

		&ret.CreatedAt,
		&ret.LastUpdate,

		&ret.GiverReview.Rating,
		&ret.GiverReview.Comment,
		&ret.TakerReview.Rating,
		&ret.TakerReview.Comment,
		&ret.Status,
	)

	return &ret, t.conn.CheckError(err)
}
