package category

import (
	"log"

	storage "github.com/alexwbaule/give-help/v2/internal/storage/category"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

//Category Object struct
type Category struct {
	storage *storage.Category
}

//New creates a new instance
func New(conn *connection.Connection) *Category {
	return &Category{
		storage: storage.New(conn),
	}
}

//Load load categories
func (c *Category) Load() ([]string, error) {

	ret, err := c.storage.Load()

	if err != nil {
		log.Printf("fail to try load categories: %s", err)
	}

	return ret, err
}

//Insert insert categories
func (c *Category) Insert(categories []string) error {
	if len(categories) < 1 {
		return nil
	}

	qtd, err := c.storage.Insert(categories)

	if qtd != int64(len(categories)) {
		log.Printf("numer of affected rows is invalid. Expected %d, received %d", qtd, len(categories))
	}

	if err != nil {
		log.Printf("fail to try insert categories [%v]: %s", categories, err)
	}

	return err
}
