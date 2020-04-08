package users

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage"

	"testing"
)

func TestUser(t *testing.T) {
	dbConfig := &storage.Config{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := storage.New(dbConfig)

	userStorage := New(conn)

	userID := "01E5DEKKFZRKEYCRN6PDXJ8GYZ"

	userData := &models.User{
		UserID:         models.ID(userID),
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
		Name:        "Usuario Da Silva",
		Reputation: &models.Reputation{
			Giver: 2.5,
			Taker: 2.5,
		},
		Tags: models.Tags([]string{"Usuário de testes", "TI", "Serviços Gerais"}),
		Location: &models.Location{
			Address: "Rua da casa do usuário, 777",
			City:    "São Paulo",
			Country: "Brasil",
			State:   "São Paulo",
			ZipCode: 99999000,
			Lat:     -23.5475,
			Long:    -46.63611,
		},
	}

	err := userStorage.Upsert(userData)

	if err != nil {
		t.Errorf("fail to try insert user data from %v - error: %s", userData, err)
	}

	u, err := userStorage.Load(userID)

	if err != nil {
		t.Errorf("fail to load user, error=%s", err)
	}

	u.Description = "Nosso querido usuário de testes unitários, agora atualizado"
	u.Contact.Phones[0].PhoneNumber = "88888-8888"

	err = userStorage.Upsert(u)

	if err != nil {
		t.Errorf("fail to try update user data from %v - error: %s", userData, err)
	}

	updated, err := userStorage.Load(userID)

	if err != nil {
		t.Errorf("fail to load user, error=%s", err)
	}

	if updated.Description == userData.Description {
		t.Errorf("fail to update user (description), error=%s", err)
	}

	if updated.Contact.Phones[0].PhoneNumber == userData.Contact.Phones[0].PhoneNumber {
		t.Errorf("fail to update user (phone number), error=%s", err)
	}
}
