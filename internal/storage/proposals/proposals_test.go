package proposals

import (
	"fmt"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage"
	"github.com/go-openapi/strfmt"

	"testing"
)

func createConn() *Proposals {
	dbConfig := &storage.Config{
		Host:   "localhost",
		User:   "postgres",
		Pass:   "example",
		DBName: "postgres",
	}

	conn := storage.New(dbConfig)

	return New(conn)
}

func getUserID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GYZ"
}

func getPrposalID() string {
	return "01E5DEKKFZRKEYCRN6PDXJ8GXX"
}

func createProposal() *models.Proposal {
	proposalID := getPrposalID()
	userID := getUserID()

	return &models.Proposal{
		ProposalID:       models.ID(proposalID),
		UserID:           models.ID(userID),
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
		Description: "Estou morrendo de fome, adoraria qualquer coisa para comer",
		Tags:        models.Tags([]string{"Alimentação"}),
	}
}

func TestUpsert(t *testing.T) {
	storage := createConn()

	proposalID := getPrposalID()
	data := createProposal()

	err := storage.Upsert(data)

	if err != nil {
		t.Errorf("fail to try insert proposal data from %v - error: %s", data, err)
	}

	loaded, err := storage.LoadFromProposal(proposalID)

	if err != nil {
		t.Errorf("fail to load proposal, error=%s", err)
	}

	if loaded.UserID != data.UserID {
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

	if loaded.TargetArea.Lat != data.TargetArea.Lat {
		t.Errorf("fail to load proposal, [TargetArea.Lat] expected: %f received: %f", data.TargetArea.Lat, loaded.TargetArea.Lat)
	}

	if loaded.TargetArea.Long != data.TargetArea.Long {
		t.Errorf("fail to load proposal, [TargetArea.Long] expected: %f received: %f", data.TargetArea.Long, loaded.TargetArea.Long)
	}

	if loaded.TargetArea.Range != data.TargetArea.Range {
		t.Errorf("fail to load proposal, [TargetArea.Range] expected: %f received: %f", data.TargetArea.Range, loaded.TargetArea.Range)
	}

	if len(loaded.TargetArea.AreaTags) != len(data.TargetArea.AreaTags) {
		t.Errorf("fail to load proposal, [TargetArea.AreaTags] expected: %s received: %s", data.TargetArea.AreaTags, loaded.TargetArea.AreaTags)
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
	return &models.Filter{
		TargetArea: &models.Area{
			Lat:   -23.5475,
			Long:  -46.6361,
			Range: 15,
		},
	}
}

func createFilterAreaTags() *models.Filter {
	return &models.Filter{
		TargetArea: &models.Area{
			AreaTags: models.Tags{"ZL"},
			Lat:      -23.5475,
			Long:     -46.6361,
			Range:    15,
		},
	}
}

func createFilterInvalidArea() *models.Filter {
	return &models.Filter{
		TargetArea: &models.Area{
			Lat:   -20.5475,
			Long:  -40.6361,
			Range: 5,
		},
	}
}

func createFilterAll() *models.Filter {
	return &models.Filter{
		Description:   "fome",
		Side:          models.SideRequest,
		ProposalTypes: []models.Type{models.TypeProduct, models.TypeService},
		TargetArea: &models.Area{
			Lat:   -23.5475,
			Long:  -46.6361,
			Range: 10,
		},
	}
}

func TestFilter(t *testing.T) {
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

	//not match
	proposals, err = storage.Find(createFilterInvalidArea())

	if err != nil {
		t.Errorf("fail to find proposals, error=%s", err)
	}

	if len(proposals) > 0 {
		t.Error("proposals area filter didn't work! (found an invalid area)")
	}
}
