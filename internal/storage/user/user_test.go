package user

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"

	"testing"
)

func createConn() *User {
	dbConfig := &common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := connection.New(dbConfig)

	return New(conn)
}

func getUserID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
}

func createUser() *models.User {
	userID := getUserID()

	giver := float64(2.6)
	taker := float64(2.5)
	zipcode := int64(11222333)
	lat := float64(-23.5475)
	lon := float64(-46.63611)

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
				&models.Phone{
					CountryCode: "+55",
					IsDefault:   false,
					PhoneNumber: "92222-2222",
					Region:      "11",
					Whatsapp:    false,
				},
				&models.Phone{
					CountryCode: "+55",
					IsDefault:   false,
					PhoneNumber: "93333-3333",
					Region:      "11",
					Whatsapp:    false,
				},
			},
		},
		Description: "Nosso querido usuário de testes unitários",
		DeviceID:    common.GetULID(),
		Name:        "Usuario Da Silva",
		Reputation: &models.Reputation{
			Giver: &giver,
			Taker: &taker,
		},
		Tags: models.Tags([]string{"Usuário de testes", "TI", "Serviços Gerais"}),
		Location: &models.Location{
			Address: "Rua da casa do usuário, 777",
			City:    "São Paulo",
			Country: "Brasil",
			State:   "São Paulo",
			ZipCode: &zipcode,
			Lat:     &lat,
			Lon:     &lon,
		},
	}
}

func insert(t *testing.T) {
	userStorage := createConn()

	userData := createUser()
	userID := string(userData.UserID)

	err := userStorage.Upsert(userData)

	if err != nil {
		t.Errorf("fail to try insert user data from %v - error: %s", userData, err)
	}

	_, err = userStorage.Load(userID)

	if err != nil {
		t.Errorf("fail to load user, error=%s", err)
	}
}

func update(t *testing.T) {
	userStorage := createConn()

	userData := createUser()
	userID := getUserID()

	userLoaded, err := userStorage.Load(userID)

	if err != nil {
		t.Errorf("fail to load user, error=%s", err)
	}

	userLoaded.Description = "Nosso querido usuário de testes unitários, agora atualizado"
	userLoaded.Contact.Phones[1].PhoneNumber = "88888-8888"
	userLoaded.Contact.Phones[1].Whatsapp = true

	err = userStorage.Upsert(userLoaded)

	if err != nil {
		t.Errorf("fail to try update user data from %v - error: %s", userData, err)
	}

	updated, err := userStorage.Load(string(userLoaded.UserID))

	if err != nil {
		t.Errorf("fail to load user, error=%s", err)
	}

	if updated.Description == userData.Description {
		t.Errorf("fail to update user (description), error=%s", err)
	}

	if updated.Contact.Phones[1].PhoneNumber == userData.Contact.Phones[1].PhoneNumber {
		t.Errorf("fail to update user (phone number), error=%s", err)
	}

	if updated.Contact.Phones[1].Whatsapp == userData.Contact.Phones[1].Whatsapp {
		t.Errorf("fail to update user (Whatsapp), error=%s", err)
	}
}

func loadAll(t *testing.T) {
	userStorage := createConn()

	users, err := userStorage.LoadAll()

	if err != nil {
		t.Errorf("fail to load user, error=%s", err)
	}

	if len(users) == 0 {
		t.Errorf("fail to load user, error=%s", err)
	}
}

func Test(t *testing.T) {
	insert(t)
	update(t)
	loadAll(t)
}
