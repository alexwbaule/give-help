package proposal

import (
	"fmt"
	"log"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	conn "github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/proposal"
	"github.com/go-openapi/strfmt"
)

//Proposal Object struct
type Proposal struct {
	storage *storage.Proposal
	config  *common.Config
}

//New creates a new instance
func New(config *common.Config) *Proposal {
	conn := conn.New(config.Db)
	return &Proposal{
		storage: storage.New(conn),
		config:  config,
	}
}

//Insert insert data
func (p *Proposal) Insert(proposal *models.Proposal) (models.ID, error) {
	if len(proposal.ProposalID) == 0 {
		proposal.ProposalID = models.ID(common.GetULID())
	}

	err := p.storage.Upsert(proposal)

	if err != nil {
		log.Printf("fail to insert new proposal [%s]: %s", proposal.ProposalID, err)
	}

	return proposal.ProposalID, err
}

func (p *Proposal) update(proposal *models.Proposal) error {
	if len(proposal.ProposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	err := p.storage.Upsert(proposal)

	if err != nil {
		log.Printf("fail to update proposal [%s]: %s", proposal.ProposalID, err)
	}

	return err
}

//LoadFromProposal load data
func (p *Proposal) Load(proposalID string) (*models.Proposal, error) {
	if len(proposalID) == 0 {
		return nil, fmt.Errorf("proposalID is empty")
	}

	ret, err := p.storage.LoadFromProposal(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
	}

	return ret, err
}

//Load load data from user
func (p *Proposal) LoadFromUser(userID string) ([]*models.Proposal, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("userID is empty")
	}

	ret, err := p.storage.LoadFromUser(userID)

	if err != nil {
		log.Printf("fail to load userID [%s] proposals: %s", userID, err)
	}

	return ret, err
}

//Find find all proposals that match with filter
func (p *Proposal) Find(filter *models.Filter) ([]*models.Proposal, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is null")
	}

	ret, err := p.storage.Find(filter)

	if err != nil {
		log.Printf("fail to load data from filter: %s", err)
	}

	return ret, err
}

//GetUserDataToShare load user data to share on proposal
func (p *Proposal) GetUserDataToShare(proposalID string) ([]models.DataToShare, error) {
	if len(proposalID) == 0 {
		return []models.DataToShare{}, fmt.Errorf("proposalID is empty")
	}

	ret, err := p.Load(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
	}

	return ret.DataToShare, err
}

//ChangeValidStatus change active field
func (p *Proposal) ChangeValidStatus(proposalID string, status bool) error {
	if len(proposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	prop, err := p.Load(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return err
	}

	prop.IsActive = status

	return p.update(prop)
}

//AddTags add proposal tags
func (p *Proposal) AddTags(proposalID string, tags []string) error {
	if len(proposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	prop, err := p.Load(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return err
	}

	prop.Tags = append(prop.Tags, tags...)

	return p.update(prop)
}

//AddImages add proposal images
func (p *Proposal) AddImages(proposalID string, images []string) error {
	if len(proposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	prop, err := p.Load(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return err
	}

	prop.Images = append(prop.Images, images...)

	return p.update(prop)
}

//ChangeValidate change proposal validate
func (p *Proposal) ChangeValidate(proposalID string, validate time.Time) error {
	if len(proposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	prop, err := p.Load(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return err
	}

	prop.ProposalValidate = strfmt.DateTime(validate)

	return p.update(prop)
}

//ChangeText change proposal title and description
func (p *Proposal) ChangeText(proposalID string, title string, description string) error {
	if len(proposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	prop, err := p.Load(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return err
	}

	prop.Title = title
	prop.Description = description

	return p.update(prop)
}
