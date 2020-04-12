package category

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func createHandler() *Category {
	c := connection.New(&common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	})

	return New(c)
}

func getServiceCategory() string {
	return "Service Test Category"
}

func TestCategories(t *testing.T) {
	service := createHandler()

	err := service.Insert([]string{getServiceCategory()})

	if err != nil {
		t.Errorf("fail to load categories: %s", err.Error())
	}

	cat, err := service.Load()

	for _, c := range cat {
		if c == getServiceCategory() {
			return
		}
	}

	t.Errorf("fail to insert categories: %s", err.Error())
}
