package bank

import (
	"context"
	"fmt"
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

//Banks Object struct
type Bank struct {
	conn *connection.Connection
}

//New creates a new instance
func New(conn *connection.Connection) *Bank {
	return &Bank{conn: conn}
}

const selectBanks = `
SELECT
	BankID,
	BankName,
	BankFullName
FROM
	Banks
ORDER BY
	BankID
`

func (b *Bank) LoadBanks() ([]*models.Bank, error) {
	ret := []*models.Bank{}

	db := b.conn.Get()

	rows, err := db.Query(selectBanks)

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			bank := &models.Bank{}
			if err = rows.Scan(
				&bank.BankID,
				&bank.BankName,
				&bank.BankFullname,
			); err == nil {
				ret = append(ret, bank)
			} else {
				return ret, b.conn.CheckError(err)
			}
		}
	}

	return ret, b.conn.CheckError(err)
}

const insertAccounts = `INSERT INTO BANK_ACCOUNT 
(
	AccountNumber,
	AccountDigit,
	AccountOwner,
	AcountDocument,
	BranchNumber,
	BranchDigit,
	BankID,
	ProposalID
) 
VALUES 
(
	$1, --AccountNumber,
	$2, --AccountDigit,
	$3, --AccountOwner,
	$4, --AcountDocument,
	$5, --BranchNumber,
	$6, --BranchDigit
	$7, --BankID,
	$8 --ProposalID
)
ON CONFLICT (Name) 
DO NOTHING;`

const removeProposalAccounts = `
DELETE FROM BANK_ACCOUNT WHERE ProposalID = $1;
`

//Insert insert categories on database
func (b *Bank) UpsertAccounts(proposalID string, accs []*models.BankAccount) error {
	db := b.conn.Get()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	//clean proposal accounts
	if _, err := db.ExecContext(ctx, removeProposalAccounts, proposalID); err != nil {
		log.Printf("fail to try clean proposal bank accounts, calling rollback: %s", err)
		tx.Rollback()
		return b.conn.CheckError(err)
	}

	for _, acc := range accs {
		result, err := db.ExecContext(
			ctx,
			insertAccounts,
			acc.AccountNumber,
			acc.AccountDigit,
			acc.AccountOwner,
			acc.AccountDocument,
			acc.BranchNumber,
			acc.BranchDigit,
			acc.BankID,
			proposalID,
		)

		if err != nil {
			log.Printf("fail to try insert new proposal bank accounts (insert fail), calling rollback: %s", err)
			tx.Rollback()
			return b.conn.CheckError(err)
		}

		aff, err := result.RowsAffected()

		if err != nil {
			log.Printf("fail to try insert new proposal bank accounts (read result), calling rollback: %s", err)
			tx.Rollback()
			return b.conn.CheckError(err)
		}

		if aff == 0 {
			log.Printf("fail to try insert new proposal bank accounts (no rows affected), calling rollback: %s", err)
			tx.Rollback()
			return fmt.Errorf("0 rows affected, check arguments!")
		}
	}

	return tx.Commit()
}

const selectAccounts = `
SELECT 
	B.BankID,
	B.BankName,
	B.BankFullName,
	A.AccountNumber,
	A.AccountDigit,
	A.AccountOwner,
	A.AcountDocument,
	A.BranchNumber,
	A.BranchDigit
FROM 
	BANK_ACCOUNTS A INNER JOIN BANKS B
		ON A.BankID = B.BankID
WHERE
	A.ProposalID = $1
ORDER BY 
	CreatedAt;
`

//LoadAccounts load bank account from a proposal
func (b *Bank) LoadAccounts(proposalId string) ([]*models.BankAccount, error) {
	ret := []*models.BankAccount{}

	db := b.conn.Get()

	rows, err := db.Query(selectAccounts, proposalId)

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			acc := &models.BankAccount{}

			if err = rows.Scan(
				&acc.BankID,
				&acc.BankName,
				&acc.BankFullname,
				&acc.AccountNumber,
				&acc.AccountDigit,
				&acc.AccountOwner,
				&acc.AccountDocument,
				&acc.BranchNumber,
				&acc.BranchDigit,
			); err == nil {
				ret = append(ret, acc)
			} else {
				return ret, b.conn.CheckError(err)
			}
		}
	}

	return ret, b.conn.CheckError(err)
}
