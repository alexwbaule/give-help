package tags

import (
	"log"

	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/tags"
)

//Tags Object struct
type Tags struct {
	storage *storage.Tags
}

//New creates a new instance
func New(conn *connection.Connection) *Tags {
	return &Tags{
		storage: storage.New(conn),
	}
}

//Load load categories
func (t *Tags) Load() ([]string, error) {

	ret, err := t.storage.Load()

	if err != nil {
		log.Printf("fail to try load tags: %s", err)
	}

	return ret, err
}

//Insert insert categories
func (t *Tags) Insert(tags []string) error {
	if len(tags) < 1 {
		return nil
	}

	qtd, err := t.storage.Insert(tags)

	if qtd != int64(len(tags)) {
		log.Printf("numer of affected rows is invalid. Expected %d, received %d", qtd, len(tags))
	}

	if err != nil {
		log.Printf("fail to try insert tags [%v]: %s", tags, err)
	}

	return err
}
