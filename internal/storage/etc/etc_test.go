package etc

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

	read, err := storage.Load()

	if err == nil && len(read) == 0 {
		t.Errorf("Invalid etc data load")
	}
}
