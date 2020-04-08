package categories

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/storage"
)

func TestCategories(t *testing.T) {
	dbConfig := &storage.Config{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := storage.New(dbConfig)

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
