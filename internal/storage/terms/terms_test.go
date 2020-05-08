package terms

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func getUserID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
}

func loadTerms(t *testing.T) {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := connection.New(dbConfig)

	storage := New(conn)

	terms, err := storage.LoadTerms()

	if err != nil {
		t.Errorf("Invalid terms load: %s", err)
	}

	if len(terms) == 0 {
		t.Errorf("no terms loaded")
	}
}

func acceptTerms(t *testing.T) {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := connection.New(dbConfig)

	storage := New(conn)

	terms, err := storage.LoadTerms()

	if err != nil {
		t.Errorf("Invalid terms load: %s", err)
	}

	if len(terms) == 0 {
		t.Errorf("no terms loaded")
	}

	for _, term := range terms {
		err = storage.Accept(string(term.TermID), getUserID())

		if err != nil {
			t.Errorf("fail to acceptd term: %s", err)
		}
	}
}

func loadUserTerms(t *testing.T) {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := connection.New(dbConfig)

	storage := New(conn)

	terms, err := storage.LoadUserAcceptedTerms(getUserID())

	if err != nil {
		t.Errorf("Invalid terms load: %s", err)
	}

	if len(terms) == 0 {
		t.Errorf("no terms loaded")
	}
}

func Test(t *testing.T) {
	loadTerms(t)
	acceptTerms(t)
	loadUserTerms(t)
}
