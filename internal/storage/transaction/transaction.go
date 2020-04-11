package transaction

import (
	"fmt"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	"github.com/lib/pq"
)

//Transaction Object struct
type Transaction struct {
	conn *connection.Connection
}

//New creates a new instance
func New(conn *connection.Connection) *Transaction {
	return &Transaction{conn: conn}
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
func (t *Transaction) Upsert(transaction *models.Transaction) error {
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

func (t *Transaction) LoadByProposalID(proposalID string) ([]*models.Transaction, error) {
	return t.load(fmt.Sprintf(selectTransaction, "ProposalID = $1"), proposalID)
}

func (t *Transaction) LoadByUserID(userID string) ([]*models.Transaction, error) {
	return t.load(fmt.Sprintf(selectTransaction, "GiverID = $1 OR TakerID = $1"), userID)
}

func (t *Transaction) load(cmd string, args ...interface{}) ([]*models.Transaction, error) {
	ret := []*models.Transaction{}

	db := t.conn.Get()
	defer db.Close()

	rows, err := db.Query(cmd, args...)

	if err != nil {
		return ret, t.conn.CheckError(err)
	}

	defer rows.Close()

	for rows.Next() {
		i := models.Transaction{
			GiverReview: &models.Review{},
			TakerReview: &models.Review{},
		}

		err := rows.Scan(
			&i.TransactionID,
			&i.ProposalID,
			&i.GiverID,
			&i.TakerID,
			&i.CreatedAt,
			&i.LastUpdate,
			&i.GiverReview.Rating,
			&i.GiverReview.Comment,
			&i.TakerReview.Rating,
			&i.TakerReview.Comment,
			&i.Status,
		)

		if err != nil {
			return ret, t.conn.CheckError(err)
		}

		ret = append(ret, &i)
	}

	return ret, t.conn.CheckError(err)
}
