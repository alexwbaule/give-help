package terms

import (
	"fmt"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

//Tags Object struct
type Terms struct {
	conn *connection.Connection
}

//New creates a new instance
func New(conn *connection.Connection) *Terms {
	return &Terms{conn: conn}
}

const insertAccept = `INSERT INTO TERMS_ACCEPTED 
(
	UserID,
	TermID
) 
VALUES 
(
	$1,
	$2
)
ON CONFLICT (TermID, UserID)
DO 
	UPDATE
		SET LastUpdate = CURRENT_TIMESTAMP;
`

//Accept insert on database an entry
func (t *Terms) Accept(termID string, userID string) error {
	db := t.conn.Get()

	result, err := db.Exec(insertAccept, userID, termID)

	if err != nil {
		return t.conn.CheckError(err)
	}

	aff, err := result.RowsAffected()

	if err != nil {
		return t.conn.CheckError(err)
	}

	if aff == 0 {
		return fmt.Errorf("no rows affected, check userID [%s] and termID [%s]", userID, termID)
	}

	return err
}

const selectTerms = `
SELECT
	TermID,
	CreatedAt,
	LastUpdate,
	Title,
	Description
FROM
	TERMS
WHERE
	IsActive = true
ORDER BY
	LastUpdate
`

//LoadTerms load Terms from database
func (t *Terms) LoadTerms() ([]*models.Term, error) {
	ret := []*models.Term{}

	db := t.conn.Get()

	rows, err := db.Query(selectTerms)

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			item := models.Term{IsActive: true}

			if err = rows.Scan(
				&item.TermID,
				&item.CreatedAt,
				&item.LastUpdate,
				&item.Title,
				&item.Description,
			); err == nil {
				ret = append(ret, &item)
			}
		}
	}

	return ret, t.conn.CheckError(err)
}

const selectUserTerms = `
SELECT
	TERMS.TermID,
	TERMS.CreatedAt,
	TERMS.LastUpdate,
	TERMS.Title,
	TERMS.Description,
	TERMS_ACCEPTED.LastUpdate AcceptedAt,
	TERMS_ACCEPTED.Accepted
FROM
TERMS 
	INNER JOIN TERMS_ACCEPTED
		ON TERMS.TermID = TERMS_ACCEPTED.TermID 
WHERE
	TERMS.IsActive = true
	AND UserId = $1
ORDER BY
	TERMS_ACCEPTED.LastUpdate
`

//LoadUserAcceptedTerms load user accepted terms from database
func (t *Terms) LoadUserAcceptedTerms(userID string) ([]*models.UserTerm, error) {
	ret := []*models.UserTerm{}

	db := t.conn.Get()

	rows, err := db.Query(selectUserTerms, userID)

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			term := models.Term{IsActive: true}
			item := models.UserTerm{}

			if err = rows.Scan(
				&term.TermID,
				&term.CreatedAt,
				&term.LastUpdate,
				&term.Title,
				&term.Description,
				&item.AcceptedAt,
				&item.IsAccepted,
			); err == nil {
				item.Term = &term
				ret = append(ret, &item)
			}
		}
	}

	return ret, t.conn.CheckError(err)
}
