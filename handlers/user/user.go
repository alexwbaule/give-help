package user

import (
	"fmt"
	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/user"
)

//User Object struct
type User struct {
	storage *storage.User
}

//New creates a new instance
func New(conn *connection.Connection) *User {
	return &User{
		storage: storage.New(conn),
	}
}

//Insert insert new data
func (u *User) Insert(user *models.User, uid string) (models.UserID, error) {

	user.UserID = models.UserID(uid)

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

//LoadAll load all data
func (u *User) LoadAll() ([]*models.User, error) {
	ret, err := u.storage.LoadAll()

	if err != nil {
		log.Printf("fail to load all users: %s", err)
	}

	return ret, err
}
