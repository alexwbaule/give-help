package category

import (
	"log"

	"github.com/alexwbaule/give-help/v2/internal/common"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/category"
	conn "github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

//Category Object struct
type Category struct {
	storage *storage.Category
	config  *common.Config
}

//New creates a new instance
func New(config *common.Config) *Category {
	conn := conn.New(config.Db)
	return &Category{
		storage: storage.New(conn),
		config:  config,
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
