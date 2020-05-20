package bank

import (
	"strings"
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func createHandler() *Bank {
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

	list, err := service.LoadBanks()

	if err != nil {
		t.Errorf("fail to load bank list: %s", err.Error())
	}

	for _, b := range list {
		if strings.Contains(strings.ToLower(b.BankFullname), "c6") {
			return
		}
	}

	t.Errorf("fail to load bank list, C6 not found")
}
