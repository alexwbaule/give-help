package user

import (
	"fmt"
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	conn "github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/user"
)

//User Object struct
type User struct {
	storage *storage.User
	config  *common.Config
}

//New creates a new instance
func New(config *common.Config) *User {
	conn := conn.New(config.Db)
	return &User{
		storage: storage.New(conn),
		config:  config,
	}
}

//Insert insert new data
func (u *User) Insert(user *models.User) (models.ID, error) {
	if len(user.UserID) == 0 {
		user.UserID = models.ID(common.GetULID())
	}

	err := u.storage.Upsert(user)

	if err != nil {
		log.Printf("fail to insert new user [%s]: %s", user.UserID, err)
	}

	return user.UserID, err
}

//Update update data
func (u *User) Update(user *models.User) error {
	if len(user.UserID) == 0 {
		return fmt.Errorf("userId is empty")
	}

	err := u.storage.Upsert(user)

	if err != nil {
		log.Printf("fail to update user [%s]: %s", user.UserID, err)
	}

	return err
}

//Load load data
func (u *User) Load(userID string) (*models.User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("userId is empty")
	}

	ret, err := u.storage.Load(userID)

	if err != nil {
		log.Printf("fail to load user [%s]: %s", userID, err)
	}

	return ret, err
}
