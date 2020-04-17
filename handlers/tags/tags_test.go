package tags

import (
	"testing"

	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func createHandler() *Tags {
	c := connection.New(&common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	})

	return New(c)
}

func getServiceTag() string {
	return "Service Test Tags"
}

func TestTags(t *testing.T) {
	service := createHandler()

	err := service.Insert([]string{getServiceTag()})

	if err != nil {
		t.Errorf("fail to load tags: %s", err.Error())
	}

	tag, err := service.Load()

	for _, c := range tag {
		if c == getServiceTag() {
			return
		}
	}

	t.Errorf("fail to insert tags: %s", err.Error())
}
