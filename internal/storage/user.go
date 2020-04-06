package storage

import "encoding/json"

func (s *Storage) UpsertUser(user *models.User) error {
	payload, err := json.Marshal(user)

	if err != nil {
		return err
	}

	_, err = s.send(payload)

	return err
}

func (s *Storage) LoadUser(userID string) (*models.User, error) {

}
