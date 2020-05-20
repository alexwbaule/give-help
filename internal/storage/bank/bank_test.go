package bank

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func Test(t *testing.T) {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := connection.New(dbConfig)

	storage := New(conn)

	read, err := storage.LoadBanks()

	if err != nil {
		t.Errorf("fail to load banks load, error: %s", err)
	}

	if len(read) == 0 {
		t.Errorf("fail to load banks load, error: 0 items load")
	}
}
