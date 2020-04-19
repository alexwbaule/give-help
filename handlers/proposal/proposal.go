package proposal

import (
	"fmt"
	"log"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/proposal"
	userStorage "github.com/alexwbaule/give-help/v2/internal/storage/user"
	"github.com/go-openapi/strfmt"
)

//Proposal Object struct
type Proposal struct {
	storage     *storage.Proposal
	userStorage *userStorage.User
}

//New creates a new instance
func New(conn *connection.Connection) *Proposal {
	return &Proposal{
		storage:     storage.New(conn),
		userStorage: userStorage.New(conn),
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

//LoadFromID load data
func (p *Proposal) LoadFromID(proposalID string) (*models.Proposal, error) {
	if len(proposalID) == 0 {
		return nil, fmt.Errorf("proposalID is empty")
	}

	ret, err := p.storage.LoadFromID(proposalID)

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

//LoadFromFilter Load all proposals that match with filter
func (p *Proposal) LoadFromFilter(filter *models.Filter) (*models.ProposalsResponse, error) {
	result, err := p.storage.Find(filter)

	if err != nil {
		log.Printf("fail to load data from filter: %s", err)
		return &models.ProposalsResponse{
			Filter: filter,
		}, err
	}

	ret := models.ProposalsResponse{
		Filter: filter,
		Result: result,
	}

	*ret.CurrentPage = filter.PageNumber
	*ret.CurrentPageSize = int64(len(result))

	return &ret, err
}

//GetUserDataToShare load user data to share on proposal
func (p *Proposal) GetUserDataToShare(proposalID string) ([]*models.DataToShareResponse, error) {
	if len(proposalID) == 0 {
		return []*models.DataToShareResponse{}, fmt.Errorf("proposalID is empty")
	}

	ret := []*models.DataToShareResponse{}

	prop, err := p.LoadFromID(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return ret, err
	}

	user, err := p.userStorage.Load(string(prop.UserID))

	if err != nil {
		log.Printf("fail to load user [%s]: %s", prop.UserID, err)
		return ret, err
	}

	for _, dts := range prop.DataToShare {
		if user.Contact != nil {
			switch dts {
			case models.DataToShareEmail:
				ret = append(ret, &models.DataToShareResponse{
					ContactType: models.DataToShareEmail,
					Contact:     user.Contact.Email,
				})
			case models.DataToShareFacebook:
				ret = append(ret, &models.DataToShareResponse{
					ContactType: models.DataToShareFacebook,
					Contact:     user.Contact.Facebook,
				})
			case models.DataToSharePhone:
				ret = append(ret, &models.DataToShareResponse{
					ContactType: models.DataToSharePhone,
					Contact:     user.Contact.Phones,
				})
			case models.DataToShareInstagram:
				ret = append(ret, &models.DataToShareResponse{
					ContactType: models.DataToShareInstagram,
					Contact:     user.Contact.Instagram,
				})
			case models.DataToShareURL:
				ret = append(ret, &models.DataToShareResponse{
					ContactType: models.DataToShareURL,
					Contact:     user.Contact.URL,
				})
			}
		}
	}

	return ret, err
}

//ChangeValidStatus change active field
func (p *Proposal) ChangeValidStatus(proposalID string, status bool) error {
	if len(proposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	prop, err := p.LoadFromID(proposalID)

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

	prop, err := p.LoadFromID(proposalID)

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

	prop, err := p.LoadFromID(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return err
	}

	prop.Images = append(prop.Images, images...)

	return p.update(prop)
}

//ChangeImages change proposal images
func (p *Proposal) ChangeImages(proposalID string, images []string) error {
	if len(proposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	prop, err := p.LoadFromID(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return err
	}

	prop.Images = images

	return p.update(prop)
}

//ChangeValidate change proposal validate
func (p *Proposal) ChangeValidate(proposalID string, validate time.Time) error {
	if len(proposalID) == 0 {
		return fmt.Errorf("proposalID is empty")
	}

	prop, err := p.LoadFromID(proposalID)

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

	prop, err := p.LoadFromID(proposalID)

	if err != nil {
		log.Printf("fail to load proposal [%s]: %s", proposalID, err)
		return err
	}

	prop.Title = title
	prop.Description = description

	return p.update(prop)
}
