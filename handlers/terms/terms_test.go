package terms

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func createHandler() *Terms {
	c := connection.New(&common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	})

	return New(c)
}

func getUserID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
}

func TestTerms(t *testing.T) {
	service := createHandler()

	trm, err := service.LoadTerms()

	if err != nil {
		t.Errorf("fail to load terms: %s", err.Error())
	}

	if len(trm) == 0 {
		t.Errorf("no terms loaded")
	}

	err = service.Accept(string(trm[0].TermID), getUserID())

	if err != nil {
		t.Errorf("fail to try accept terms: %s", err.Error())
	}

	userTerms, err := service.LoadUserAcceptedTerms(getUserID())

	if err != nil {
		t.Errorf("fail to load user terms: %s", err.Error())
	}

	if len(userTerms) == 0 {
		t.Errorf("no user term loaded")
	}
}
