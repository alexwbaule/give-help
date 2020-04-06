package storage

import "encoding/json"

func (s *Storage) InsertTransaction(transaction *models.Transaction) error {
	payload, err := json.Marshal(transaction)

	if err != nil {
		return err
	}

	_, err = s.send(payload)

	return err
}

func (s *Storage) InsertGiverReview(transactionID string, review *models.Review) error {

}

func (s *Storage) InsertTakerReview(transactionID string, review *models.Review) error {

}

func (s *Storage) UpdateTransactionStatus(transactionID string, status models.TransactionStatus) error {
	return nil
}
