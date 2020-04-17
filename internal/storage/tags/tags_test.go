package tags

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func TestTags(t *testing.T) {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := connection.New(dbConfig)

	storage := New(conn)

	tags := []string{"Construção Cívil", "Pets", "Mecânica", "Alimentação"}

	qtd, err := storage.Insert(tags)

	if int64(len(tags)) < qtd {
		t.Errorf("Invalid length of tags inserted, expected %d received %d", len(tags), qtd)
	}

	read, err := storage.Load()

	if err == nil && len(read) == 0 {
		t.Errorf("Invalid tags load")
	}
}
