package storage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/alexwbaule/give-help/v2/generated/models"
)

type Storage struct {
	postURL string
}

func New(postUrl string) (*Storage, error) {
	return &Storage{postURL: postUrl}, nil
}

func (s *Storage) send(payload []byte) ([]byte, error) {
	resp, err := http.Post(s.postURL, "application/json", bytes.NewBuffer(payload))
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
	return nil
}

func (s *Storage) InsertCategories(categories []string) error {
	return nil
}

func (s *Storage) GetCategories() ([]string, error) {
	return []string{""}, nil
}
