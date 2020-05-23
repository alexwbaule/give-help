package etc

import (
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/etc"
)

//Etc Object struct
type Etc struct {
	storage *storage.Etc
}

//New creates a new instance
func New(conn *connection.Connection) *Etc {
	return &Etc{
		storage: storage.New(conn),
	}
}

//Load load etc key value data
func (e *Etc) Load() (models.Etc, error) {
	ret, err := e.storage.Load()

	if err != nil {
		log.Printf("fail to try load etc key value data: %s", err)
	}

	return ret, err
}
