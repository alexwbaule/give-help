package proposal

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	cacheConnection "github.com/alexwbaule/give-help/v2/internal/cache/connection"
	"github.com/alexwbaule/give-help/v2/internal/common"
	dbConnection "github.com/alexwbaule/give-help/v2/internal/storage/connection"
	"github.com/go-openapi/strfmt"
)

func createHandler() *Proposal {
	dbConn := dbConnection.New(&common.DbConfig{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	})

	cacheConn, err := cacheConnection.New(&common.CacheConfig{
		Addresses: []string{"http://localhost:9200"},
	})

	if err != nil {
		panic(err)
	}

	return New(dbConn, cacheConn)
}

func getUserID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8UUU"
}

func getPrposalID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8BBB"
}

func createProposal() *models.Proposal {
	proposalID := getPrposalID()
	userID := getUserID()

	lat := float64(-23.5475)
	long := float64(-46.6361)
	estimatedValue := float64(99.3)

	return &models.Proposal{
		ProposalID:       models.ID(proposalID),
		UserID:           models.UserID(userID),
		IsActive:         true,
		ProposalType:     models.TypeProduct,
		Side:             models.SideLocalBusiness,
		ProposalValidate: strfmt.DateTime(time.Time{}.AddDate(2020, 5, 8)),
		TargetArea: &models.Location{
			AreaTags: models.Tags([]string{"ZL", "Penha", "Zona Leste"}),
			Lat:      &lat,
			Long:     &long,
			Range:    5,
			City:     "Porto Alegre",
			State:    "RS",
			Country:  "Brasil",
		},
		Title:          "Quero comer",
		Description:    "Estou morrendo de fome, adoraria qualquer coisa para comer",
		Tags:           models.Tags([]string{"Alimentação", "Comercio"}),
		Images:         []string{`http://my-domain.com/image1.jpg`, `http://my-domain.com/image2.jpg`, `http://my-domain.com/image3.jpg`},
		EstimatedValue: &estimatedValue,
		ExposeUserData: true,
		DataToShare:    []models.DataToShare{models.DataToSharePhone, models.DataToShareEmail, models.DataToShareFacebook, models.DataToShareInstagram, models.DataToShareURL, models.DataToShareBankAccount},
		BankAccounts: []*models.BankAccount{
			{
				BankID:          336,
				AccountNumber:   "196",
				AccountDigit:    "1",
				BranchNumber:    "00001",
				AccountDocument: "987.413.858-02",
				AccountOwner:    "Todos Juntos Ajudando",
			},
		},
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

func testInsert(t *testing.T) {
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

func testLoadFromID(t *testing.T) {
	service, prop := prepare(t)

	loaded, err := service.LoadFromID(getPrposalID())

	if err != nil {
		t.Errorf("fail to try LoadFromID proposal data from %v - error: %s", prop, err.Error())
	}

	if prop.ProposalID != loaded.ProposalID {
		t.Errorf("fail to try LoadFromID proposal data from %s", getUserID())
	}
}

func testLoadFromIDFromUser(t *testing.T) {
	service, prop := prepare(t)

	props, err := service.LoadFromUser(string(prop.UserID))

	if err != nil {
		t.Errorf("fail to try LoadFromID proposal data from %v - error: %s", props, err.Error())
	}

	if len(props) == 0 {
		t.Errorf("fail to try LoadFromID proposal data from user %s", getUserID())
	}
}

func testDTS(t *testing.T) {
	service, prop := prepare(t)

	dts, err := service.GetUserDataToShare(getPrposalID())

	if err != nil {
		t.Errorf("fail to try LoadFromID proposal data to share from %v - error: %s", prop, err.Error())
	}

	if len(dts) != len(prop.DataToShare) {
		t.Errorf("fail to try LoadFromID proposal data to share from (invalid data) %v - error: %s", prop, err.Error())
	}
}

func testChangeValidate(t *testing.T) {
	service, prop := prepare(t)

	newValidate := time.Time{}.AddDate(2020, 6, 8)
	if err := service.ChangeValidate(getPrposalID(), newValidate); err == nil {
		if loaded, err := service.LoadFromID(getPrposalID()); err == nil {
			if newValidate.Unix() != time.Time(loaded.ProposalValidate).Unix() {
				t.Errorf("invalid loaded value - propoosal (validate) expected: %s received: %s", newValidate, loaded.ProposalValidate)
			}
		} else {
			t.Errorf("fail to try LoadFromID updated proposal (validate) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (validate) from %v - error: %s", prop, err.Error())
	}
}

func testChangeValidStatus(t *testing.T) {
	service, prop := prepare(t)

	newStatus := false
	if err := service.ChangeValidStatus(getPrposalID(), newStatus); err == nil {
		if loaded, err := service.LoadFromID(getPrposalID()); err == nil {
			if newStatus != loaded.IsActive {
				t.Errorf("invalid loaded value - propoosal (IsActive) expected: %v received: %v", newStatus, loaded.IsActive)
			}
		} else {
			t.Errorf("fail to try LoadFromID updated proposal (IsActive) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (IsActive) from %v - error: %s", prop, err.Error())
	}
}

func testAddTags(t *testing.T) {
	service, prop := prepare(t)

	newTag := common.NormalizeTagArray([]string{"TestingService"})

	if err := service.AddTags(getPrposalID(), newTag); err == nil {
		if loaded, err := service.LoadFromID(getPrposalID()); err == nil {

			found := 0

			for _, t := range loaded.Tags {
				for _, nt := range newTag {
					if nt == t {
						found++
					}
				}
			}

			if found != len(newTag) {
				t.Errorf("invalid loaded value - propoosal (AddTags) tag not found!")
			}

		} else {
			t.Errorf("fail to try LoadFromID updated proposal (AddTags) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (AddTags) from %v - error: %s", prop, err.Error())
	}
}

func testAddImages(t *testing.T) {
	service, prop := prepare(t)

	newImage := "http://my-domain-test.com/image-test-1.jpg"

	if err := service.AddImages(getPrposalID(), []string{newImage}); err == nil {
		if loaded, err := service.LoadFromID(getPrposalID()); err == nil {

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
			t.Errorf("fail to try LoadFromID updated proposal (AddImages) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (AddImages) from %v - error: %s", prop, err.Error())
	}
}

func testChangeText(t *testing.T) {
	service, prop := prepare(t)

	newTitle := "Estou com fome e testando o código"
	newDesc := "Sim, dá fome testar tanto código assim, e segundo meu amigo Danilo é muito importante testar tudo direitinho, nunca vou esquece disso, já me salvou a pele várias vezes! Fica aqui a minha dica"

	if err := service.ChangeText(getPrposalID(), newTitle, newDesc); err == nil {
		if loaded, err := service.LoadFromID(getPrposalID()); err == nil {

			if newTitle != loaded.Title {
				t.Errorf("invalid loaded value - propoosal (ChangeText) expected: %s received: %s", newTitle, loaded.Title)
			}

			if newDesc != loaded.Description {
				t.Errorf("invalid loaded value - propoosal (ChangeText) expected: %s received: %s", newTitle, loaded.Title)
			}

		} else {
			t.Errorf("fail to try LoadFromID updated proposal (ChangeText) from %v - error: %s", getPrposalID(), err.Error())
		}
	} else {
		t.Errorf("fail to try update proposal (ChangeText) from %v - error: %s", prop, err.Error())
	}
}

func testFind(t *testing.T) {
	filter := &models.Filter{
		Description: "fome",
	}

	service, prop := prepare(t)

	result, err := service.LoadFromFilter(filter)

	if err != nil {
		t.Errorf("fail to try LoadFromID proposal data from %v - error: %s", prop, err.Error())
	}

	if len(result.Result) == 0 {
		t.Errorf("fail to try find data with filters - error: %s", err.Error())
	}

	if *result.CurrentPageSize < 1 {
		t.Errorf("no proposals return")
	}

	result, err = service.LoadFromFilter(nil)

	if err != nil {
		t.Errorf("fail to try LoadFromID proposal data from %v - error: %s", prop, err.Error())
	}

	if len(result.Result) == 0 {
		t.Errorf("fail to try find data with filters - error: %s", err.Error())
	}

	if *result.CurrentPageSize < 1 {
		t.Errorf("no proposals return")
	}
}

func testFindLocalBusiness(t *testing.T) {
	filter := &models.Filter{
		Side: models.SideLocalBusiness,
	}

	service, prop := prepare(t)

	result, err := service.LoadFromFilter(filter)

	if err != nil {
		t.Errorf("fail to try LoadFromID proposal data from %v - error: %s", prop, err.Error())
	}

	if len(result.Result) == 0 {
		t.Errorf("fail to try find data with filters - error: %s", err.Error())
	}

	if *result.CurrentPageSize < 1 {
		t.Errorf("no proposals return")
	}
}

func testFindOmini(t *testing.T) {
	filter := &models.Filter{
		Description: "Comercio",
	}

	service, prop := prepare(t)

	result, err := service.LoadFromFilter(filter)

	if err != nil {
		t.Errorf("fail to try LoadFromID proposal data from %v - error: %s", prop, err.Error())
	}

	if len(result.Result) == 0 {
		t.Errorf("fail to try find data with filters - error: %s", err.Error())
	}

	if *result.CurrentPageSize < 1 {
		t.Errorf("no proposals return")
	}
}

func testComplaint(t *testing.T) {
	service, prop := prepare(t)

	complaint := &models.Complaint{
		Comment:    "Não curti esse comentário via serviço",
		Complainer: "Handler do testador de sistemas",
		ProposalID: prop.ProposalID,
	}

	err := service.InsertComplaint(complaint)

	if err != nil {
		t.Errorf("fail to try insert a complaint: %s", err.Error())
	}
}

func testReindex(t *testing.T) {
	conn := createHandler()
	conn.Reindex()
}

func Test(t *testing.T) {
	testInsert(t)
	testLoadFromID(t)
	testLoadFromIDFromUser(t)
	testDTS(t)
	testChangeValidate(t)
	testChangeValidStatus(t)
	testAddTags(t)
	testAddImages(t)
	testChangeText(t)
	testComplaint(t)
	testFind(t)
	testFindLocalBusiness(t)
	testFindOmini(t)
	testReindex(t)
}

func TestX(t *testing.T) {
	testLoadFromID(t)
	//testFindOmini(t)
}
