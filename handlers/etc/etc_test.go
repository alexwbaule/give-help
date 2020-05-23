package etc

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func createHandler() *Etc {
	c := connection.New(&common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	})

	return New(c)
}

func Test(t *testing.T) {
	service := createHandler()

	list, err := service.Load()

	if err != nil {
		t.Errorf("fail to load etc key value list: %s", err.Error())
	}

	if len(list) == 0 {
		t.Errorf("etc list is empty")
	}
}
