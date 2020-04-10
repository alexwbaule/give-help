package category

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
)

func createHandler() *Category {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	config := &common.Config{
		Db: dbConfig,
	}

	return New(config)
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
