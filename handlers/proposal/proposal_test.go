package proposal

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	"github.com/go-openapi/strfmt"
)

func createHandler() *Proposal {
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

func getPrposalID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8PPP"
}

func createProposal() *models.Proposal {
	proposalID := getPrposalID()
	userID := getUserID()

	return &models.Proposal{
		ProposalID:       models.ID(proposalID),
		UserID:           models.UserID(userID),
		IsActive:         true,
		ProposalType:     models.TypeProduct,
		Side:             models.SideRequest,
		ProposalValidate: strfmt.DateTime(time.Time{}.AddDate(2020, 5, 8)),
		TargetArea: &models.Area{
			AreaTags: models.Tags([]string{"ZL", "Penha", "Zona Leste"}),
			Lat:      -23.5475,
			Long:     -46.6361,
			Range:    5,
		},
		Title:          "Quero comer",
		Description:    "Estou morrendo de fome, adoraria qualquer coisa para comer",
		Tags:           models.Tags([]string{"Alimentação"}),
		Images:         []string{`http://my-domain.com/image1.jpg`, `http://my-domain.com/image2.jpg`, `http://my-domain.com/image3.jpg`},
		EstimatedValue: float64(50),
		ExposeUserData: true,
		DataToShare:    []models.DataToShare{models.DataToSharePhone, models.DataToShareEmail, models.DataToShareFacebook, models.DataToShareInstagram, models.DataToShareURL},
	}
}

func prepare(t *testing.T) (*Proposal, *models.Proposal) {
	prop := createProposal()

	service := createHandler()

	id, err := service.Insert(prop)

	if err != nil {
		t.Errorf("fail to insert proposal data from %v - error: %s", prop, err.Error())
	}

	if len(id) == 0 {
		t.Errorf("fail to try insert proposal data from %v - error: %s", prop, fmt.Errorf("empty user id on return"))
	}

	return service, prop
}

func TestInsert(t *testing.T) {
	prop := createProposal()

	service := createHandler()

	id, err := service.Insert(prop)

	if err != nil {
		t.Errorf("fail to insert proposal data from %v - error: %s", prop, err.Error())
	}

	if len(id) == 0 {
		t.Errorf("fail to try insert proposal data from %v - error: %s", prop, fmt.Errorf("empty user id on return"))
	}
}

func TestLoad(t *testing.T) {
	service, prop := prepare(t)

	loaded, err := service.Load(getPrposalID())

	if err != nil {
		t.Errorf("fail to try load proposal data from %v - error: %s", prop, err.Error())
	}

	if prop.ProposalID != loaded.ProposalID {
		t.Errorf("fail to try load proposal data from %s", getUserID())
	}
}

func TestLoadFromUser(t *testing.T) {
	service, prop := prepare(t)

	props, err := service.LoadFromUser(string(prop.UserID))

	if err != nil {
		t.Errorf("fail to try load proposal data from %v - error: %s", props, err.Error())
	}

	if len(props) == 0 {
		t.Errorf("fail to try load proposal data from user %s", getUserID())
	}
}

func TestDTS(t *testing.T) {
	service, prop := prepare(t)

	dts, err := service.GetUserDataToShare(getPrposalID())

	if err != nil {
		t.Errorf("fail to try load proposal data to share from %v - error: %s", prop, err.Error())
	}

	if len(dts) != len(prop.DataToShare) {
		t.Errorf("fail to try load proposal data to share from (invalid data) %v - error: %s", prop, err.Error())
	}
}

func TestChangeValidate(t *testing.T) {
	service, prop := prepare(t)

	newValidate := time.Time{}.AddDate(2020, 6, 8)
	if err := service.ChangeValidate(getPrposalID(), newValidate); err == nil {
		if loaded, err := service.Load(getPrposalID()); err == nil {
			if newValidate.Unix() != time.Time(loaded.ProposalValidate).Unix() {
				t.Errorf("invalid loaded value - propoosal (validate) expected: %s received: %s", newValidate, loaded.ProposalValidate)
			}
		} else {
			t.Errorf("fail to try load updated proposal (validate) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (validate) from %v - error: %s", prop, err.Error())
	}
}

func TestChangeValidStatus(t *testing.T) {
	service, prop := prepare(t)

	newStatus := false
	if err := service.ChangeValidStatus(getPrposalID(), newStatus); err == nil {
		if loaded, err := service.Load(getPrposalID()); err == nil {
			if newStatus != loaded.IsActive {
				t.Errorf("invalid loaded value - propoosal (IsActive) expected: %v received: %v", newStatus, loaded.IsActive)
			}
		} else {
			t.Errorf("fail to try load updated proposal (IsActive) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (IsActive) from %v - error: %s", prop, err.Error())
	}
}

func TestAddTags(t *testing.T) {
	service, prop := prepare(t)

	newTag := "TestingService"

	if err := service.AddTags(getPrposalID(), []string{newTag}); err == nil {
		if loaded, err := service.Load(getPrposalID()); err == nil {

			found := false

			for _, t := range loaded.Tags {
				if t == newTag {
					found = true
				}
			}

			if !found {
				t.Errorf("invalid loaded value - propoosal (AddTags) tag not found!")
			}

		} else {
			t.Errorf("fail to try load updated proposal (AddTags) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (AddTags) from %v - error: %s", prop, err.Error())
	}
}

func TestAddImages(t *testing.T) {
	service, prop := prepare(t)

	newImage := "http://my-domain-test.com/image-test-1.jpg"

	if err := service.AddImages(getPrposalID(), []string{newImage}); err == nil {
		if loaded, err := service.Load(getPrposalID()); err == nil {

			found := false

			for _, t := range loaded.Images {
				if t == newImage {
					found = true
				}
			}

			if !found {
				t.Errorf("invalid loaded value - propoosal (AddImages) image not found!")
			}

		} else {
			t.Errorf("fail to try load updated proposal (AddImages) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (AddImages) from %v - error: %s", prop, err.Error())
	}
}

func TestChangeText(t *testing.T) {
	service, prop := prepare(t)

	newTitle := "Estou com fome e testando o código"
	newDesc := "Sim, dá fome testar tanto código assim, e segundo meu amigo Danilo é muito importante testar tudo direitinho, nunca vou esquece disso, já me salvou a pele várias vezes! Fica aqui a minha dica"

	if err := service.ChangeText(getPrposalID(), newTitle, newDesc); err == nil {
		if loaded, err := service.Load(getPrposalID()); err == nil {

			if newTitle != loaded.Title {
				t.Errorf("invalid loaded value - propoosal (ChangeText) expected: %s received: %s", newTitle, loaded.Title)
			}

			if newDesc != loaded.Description {
				t.Errorf("invalid loaded value - propoosal (ChangeText) expected: %s received: %s", newTitle, loaded.Title)
			}

		} else {
			t.Errorf("fail to try load updated proposal (ChangeText) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (ChangeText) from %v - error: %s", prop, err.Error())
	}
}

func TestFind(t *testing.T) {
	filter := &models.Filter{
		Description: "fome",
	}

	service, prop := prepare(t)

	loaded, err := service.Find(filter)

	if err != nil {
		t.Errorf("fail to try load proposal data from %v - error: %s", prop, err.Error())
	}

	if len(loaded) == 0 {
		t.Errorf("fail to try find data with filters - error: %s", err.Error())
	}
}
