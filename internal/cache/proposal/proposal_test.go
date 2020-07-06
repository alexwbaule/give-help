package proposal

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/cache/connection"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/go-openapi/strfmt"

	"testing"
)

func createConn() *Proposal {
	config := &common.CacheConfig{
		Addresses: []string{"http://localhost:9200"},
	}

	conn, err := connection.New(config)

	if err != nil {
		panic(err)
	}

	return New(conn)
}

func getUserID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
}

func getProposalID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GXX"
}

func createProposal() *models.Proposal {
	proposalID := getProposalID()
	userID := getUserID()

	estimatedValue := float64(72.6)
	lat := float64(-23.5475)
	lon := float64(-46.6361)

	return &models.Proposal{
		ProposalID:       models.ID(proposalID),
		UserID:           models.UserID(userID),
		IsActive:         true,
		ProposalType:     models.TypeProduct,
		Side:             models.SideRequest,
		ProposalValidate: strfmt.DateTime(time.Time{}.AddDate(2020, 5, 8)),
		TargetArea: &models.Location{
			AreaTags: models.Tags([]string{"ZL", "Penha", "Zona Leste"}),
			Lat:      &lat,
			Lon:      &lon,
			Distance: 5,
			City:     "São Paulo",
			State:    "SP",
			Country:  "Brazil il il il il",
		},
		Title:          "Quero comer",
		Description:    "Estou morrendo de fome, adoraria qualquer coisa para comer",
		Tags:           models.Tags([]string{"Alimentação", "comercio"}),
		Images:         []string{`http://my-domain.com/image1.jpg`, `http://my-domain.com/image2.jpg`, `http://my-domain.com/image3.jpg`},
		EstimatedValue: &estimatedValue,
		ExposeUserData: true,
		DataToShare:    []models.DataToShare{models.DataToSharePhone, models.DataToShareEmail, models.DataToShareFacebook, models.DataToShareInstagram, models.DataToShareURL},
	}
}

func createFilterDescription() *models.Filter {
	return &models.Filter{
		Description: "fome",
	}
}

func createFilterSide() *models.Filter {
	return &models.Filter{
		Side: models.SideRequest,
	}
}

func createFilterType() *models.Filter {
	return &models.Filter{
		ProposalTypes: []models.Type{models.TypeProduct, models.TypeService},
	}
}

func createFilterArea() *models.Filter {
	lat := float64(-23.5475)
	lon := float64(-46.6361)

	return &models.Filter{
		TargetArea: &models.Location{
			Lat:      &lat,
			Lon:      &lon,
			Distance: 15,
		},
	}
}

func createFilterAreaTags() *models.Filter {
	lat := float64(-23.5475)
	lon := float64(-46.6361)

	return &models.Filter{
		TargetArea: &models.Location{
			AreaTags: models.Tags{"ZL"},
			Lat:      &lat,
			Lon:      &lon,
			Distance: 15,
		},
	}
}

func createFilterInvalidArea() *models.Filter {
	lat := float64(-20.5475)
	lon := float64(-40.6361)

	return &models.Filter{
		TargetArea: &models.Location{
			Lat:      &lat,
			Lon:      &lon,
			Distance: 5,
		},
	}
}

func createFilterAll() *models.Filter {
	lat := float64(-23.5475)
	lon := float64(-46.6361)

	return &models.Filter{
		Description:   "fome",
		Side:          models.SideRequest,
		ProposalTypes: []models.Type{models.TypeProduct, models.TypeService},
		TargetArea: &models.Location{
			Lat:      &lat,
			Lon:      &lon,
			Distance: 10,
		},
	}
}

func createFilterActive() *models.Filter {
	return &models.Filter{
		IncludeInactive: false,
	}
}

func createFilterNotActive() *models.Filter {
	return &models.Filter{
		IncludeInactive: true,
	}
}

func upsert(t *testing.T) {
	storage := createConn()

	proposalID := getProposalID()
	data := createProposal()

	err := storage.Upsert(data)

	if err != nil {
		t.Errorf("fail to try insert proposal data from %v - error: %s", data, err)
	}

	loaded, err := storage.LoadFromID(proposalID)

	if err != nil {
		t.Errorf("fail to load proposal, error=%s", err)
	}

	if strings.TrimSpace(string(loaded.UserID)) != strings.TrimSpace(string(data.UserID)) {
		t.Errorf("fail to load proposal, [UserID] expected: %s received: %s", data.UserID, loaded.UserID)
	}

	if loaded.Side != data.Side {
		t.Errorf("fail to load proposal, [Side] expected: %s received: %s", data.Side, loaded.Side)
	}

	if loaded.ProposalType != data.ProposalType {
		t.Errorf("fail to load proposal, [ProposalType] expected: %s received: %s", data.ProposalType, loaded.ProposalType)
	}

	if fmt.Sprint(loaded.ProposalValidate) != fmt.Sprint(data.ProposalValidate) {
		t.Errorf("fail to load proposal, [ProposalValidate] expected: %s received: %s", data.ProposalValidate, loaded.ProposalValidate)
	}

	if loaded.IsActive != data.IsActive {
		t.Errorf("fail to load proposal, [IsActive] expected: %v received: %v", data.IsActive, loaded.IsActive)
	}

	if loaded.Description != data.Description {
		t.Errorf("fail to load proposal, [Description] expected: %s received: %s", data.Description, loaded.Description)
	}

	if len(loaded.Tags) != len(data.Tags) {
		t.Errorf("fail to load proposal, [Tags] expected: %v received: %v", data.Tags, loaded.Tags)
	}

	if *loaded.TargetArea.Lat != *data.TargetArea.Lat {
		t.Errorf("fail to load proposal, [TargetArea.Lat] expected: %f received: %f", *data.TargetArea.Lat, *loaded.TargetArea.Lat)
	}

	if *loaded.TargetArea.Lon != *data.TargetArea.Lon {
		t.Errorf("fail to load proposal, [TargetArea.Lon] expected: %f received: %f", *data.TargetArea.Lon, *loaded.TargetArea.Lon)
	}

	if loaded.TargetArea.Distance != data.TargetArea.Distance {
		t.Errorf("fail to load proposal, [TargetArea.Distance] expected: %f received: %f", data.TargetArea.Range, loaded.TargetArea.Distance)
	}

	if len(loaded.TargetArea.AreaTags) != len(data.TargetArea.AreaTags) {
		t.Errorf("fail to load proposal, [TargetArea.AreaTags] expected: %s received: %s", data.TargetArea.AreaTags, loaded.TargetArea.AreaTags)
	}
}

func filter(t *testing.T) {
	storage := createConn()

	//match
	proposals, err := storage.Find(createFilterAll())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) == 0 {
		t.Error("proposals all filter didn't work!")
	}

	proposals, err = storage.Find(createFilterDescription())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) == 0 {
		t.Error("proposals description filter didn't work!")
	}

	proposals, err = storage.Find(createFilterArea())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) == 0 {
		t.Error("proposals area filter didn't work!")
	}

	proposals, err = storage.Find(createFilterSide())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) == 0 {
		t.Error("proposals side filter didn't work!")
	}

	proposals, err = storage.Find(createFilterType())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) == 0 {
		t.Error("proposals type filter didn't work!")
	}

	//active
	proposals, err = storage.Find(createFilterActive())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) == 0 {
		t.Error("proposals active status filter didn't work!")
	}

	//not match
	proposals, err = storage.Find(createFilterInvalidArea())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) > 0 {
		t.Error("proposals area filter didn't work! (found an invalid area)")
	}

	//active
	proposals, err = storage.Find(createFilterNotActive())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) == 0 {
		t.Error("proposals active status filter didn't work! (found an invalid active status)")
	}
}

func Test(t *testing.T) {
	upsert(t)
	filter(t)
}
