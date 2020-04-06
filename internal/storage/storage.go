package storage

import (
	"bytes"
	"io/ioutil"
	"net/http"
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
