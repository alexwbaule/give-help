package terms

import (
	"fmt"
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/terms"
)

//Tags Object struct
type Terms struct {
	storage *storage.Terms
}

//New creates a new instance
func New(conn *connection.Connection) *Terms {
	return &Terms{
		storage: storage.New(conn),
	}
}

//Accept store term accepted from user
func (t *Terms) Accept(termID string, userID string) error {
	if len(userID) == 0 {
		return fmt.Errorf("userId is empty")
	}

	err := t.storage.Accept(termID, userID)

	if err != nil {
		log.Printf("fail to try insert term accept from user [%v]: %s", userID, err)
	}

	return err
}

//LoadTerms load terms
func (t *Terms) LoadTerms() ([]*models.Term, error) {
	ret, err := t.storage.LoadTerms()

	if err != nil {
		log.Printf("fail to try load terms: %s", err)
	}

	return ret, err
}

//LoadUserAcceptedTerms load user terms
func (t *Terms) LoadUserAcceptedTerms(userID string) ([]*models.UserTerm, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("userId is empty")
	}

	ret, err := t.storage.LoadUserAcceptedTerms(userID)

	if err != nil {
		log.Printf("fail to try load user terms: %s", err)
	}

	return ret, err
}
