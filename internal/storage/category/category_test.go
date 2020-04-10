package category

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func TestCategories(t *testing.T) {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := connection.New(dbConfig)

	catStorage := New(conn)

	cats := []string{"Construção Cívil", "Pets", "Mecânica", "Alimentação"}

	qtd, err := catStorage.Insert(cats)

	if int64(len(cats)) < qtd {
		t.Errorf("Invalid length of categories inserted, expected %d received %d", len(cats), qtd)
	}

	read, err := catStorage.Load()

	if err == nil && len(read) == 0 {
		t.Errorf("Invalid categories load")
	}
}
