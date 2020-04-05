package storage

import (
	"bytes"
	"net/http"
	"github.com/alexwbaule/give-help/v2/generated/models"
	"io/util"
	"encoding/json"
)

type Storage struct {
	postURL string
}

func New(postUrl string) (*Storage, error) {
	return &Storage{postUrl: postUrl}, nil
}

func (s *Storage) send(url string, payload []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return []byte{}, nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, nil
	}

	return body, err
}


func (s *Storage) UpsertUser(user *models.User) error {
	payload, err := json.Marshal(user)

	if err != nil {
		return err
	}

	_, err = s.send(payload)

	return err
}

func (s *Storage) InsertProposal(proposal *models.Proposal) error {
	payload, err := json.Marshal(proposal)

	if err != nil {
		return err
	}

	_, err = s.send(payload)

	return err
}

func (s *Storage) InsertTransaction(transaction *models.Transaction) error {
	payload, err := json.Marshal(transaction)

	if err != nil {
		return err
	}

	_, err = s.send(payload)

	return err
}

func (s *Storage) UpdateTransactionStatus(transactionID string, status models.TransactionStatus) error {

}

func (s *Storage) InsertCategories(categories []string) error {

}

func (s *Storage) GetCategories() []string, error {

}