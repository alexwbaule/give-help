package user

import (
	"fmt"
	"testing"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

func createHandler() *User {
	c := connection.New(&common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	})

	return New(c)
}

func getUserID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8UUU"
}

func createUser() *models.User {
	userID := getUserID()

	lat := float64(-23.5475)
	long := float64(-46.6361)
	zipCode := int64(99888777)
	rating := float64(3.2)

	return &models.User{
		UserID:         models.UserID(userID),
		AllowShareData: true,
		Contact: &models.Contact{
			Email:     "usuario@email.com",
			Facebook:  "usuario@facebook.com",
			Google:    "usuario@google.com",
			Instagram: "usuario@instagram.com",
			URL:       "usuario.com.br",
			Phones: []*models.Phone{
				&models.Phone{
					CountryCode: "+55",
					IsDefault:   true,
					PhoneNumber: "9999-9999",
					Region:      "11",
					Whatsapp:    true,
				},
				&models.Phone{
					CountryCode: "+55",
					IsDefault:   false,
					PhoneNumber: "1111-1111",
					Region:      "11",
					Whatsapp:    false,
				},
			},
		},
		Description: "Nosso querido usuário de testes unitários",
		DeviceID:    common.GetULID(),
		Name:        "José Insertido Pelo Serviço",
		Reputation: &models.Reputation{
			Giver: &rating,
			Taker: &rating,
		},
		Tags: models.Tags([]string{"Usuário de testes", "TI", "Serviços Gerais"}),
		Location: &models.Location{
			Address: "Rua da casa do usuário, 777",
			City:    "São Paulo",
			Country: "Brasil",
			State:   "São Paulo",
			ZipCode: &zipCode,
			Lat:     &lat,
			Long:    &long,
		},
	}
}

func Test(t *testing.T) {
	testUserInsert(t)
	testUserUpdate(t)
	testUserLoad(t)
}

func testUserInsert(t *testing.T) {
	user := createUser()

	service := createHandler()

	id, err := service.Insert(user, string(user.UserID))

	if err != nil {
		t.Errorf("fail to try insert user data from %v - error: %s", user, err.Error())
	}

	if len(id) == 0 {
		t.Errorf("fail to try insert user data from %v - error: %s", user, fmt.Errorf("empty user id on return"))
	}
}

func testUserUpdate(t *testing.T) {
	user := createUser()

	service := createHandler()

	user.Name = "Jose Alterado Pelo Serviço"

	err := service.Update(user)

	if err != nil {
		t.Errorf("fail to try update user data from %v - error: %s", user, err.Error())
	}
}

func testUserLoad(t *testing.T) {
	service := createHandler()

	user, err := service.Load(getUserID())

	if err != nil {
		t.Errorf("fail to try load user data from %v - error: %s", user, err.Error())
	}

	if user.UserID != models.UserID(getUserID()) {
		t.Errorf("fail to try load user data from %v", getUserID())
	}
}
