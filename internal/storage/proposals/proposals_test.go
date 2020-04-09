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
			Long:     -46.63611,
			Range:    25,
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
		t.Errorf("fail to try insert prposal data from %v - error: %s", data, err)
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
