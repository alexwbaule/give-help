package proposal

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	cacheConnection "github.com/alexwbaule/give-help/v2/internal/cache/connection"
	cache "github.com/alexwbaule/give-help/v2/internal/cache/proposal"
	"github.com/alexwbaule/give-help/v2/internal/common"
	dbConnection "github.com/alexwbaule/give-help/v2/internal/storage/connection"
	storage "github.com/alexwbaule/give-help/v2/internal/storage/proposal"
	tagsStorage "github.com/alexwbaule/give-help/v2/internal/storage/tags"
	userStorage "github.com/alexwbaule/give-help/v2/internal/storage/user"
	"github.com/go-openapi/strfmt"
)

//Proposal Object struct
type Proposal struct {
	cache   *cache.Proposal
	storage *storage.Proposal
	user    *userStorage.User
	tags    *tagsStorage.Tags
}

//New creates a new instance
func New(dbConn *dbConnection.Connection, cacheConn *cacheConnection.Connection) *Proposal {
	return &Proposal{
		storage: storage.New(dbConn),
		user:    userStorage.New(dbConn),
		tags:    tagsStorage.New(dbConn),
		cache:   cache.New(cacheConn),
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
		return proposal.ProposalID, err
	}

	_, err = p.tags.Insert(proposal.Tags)

	if err != nil {
		log.Printf("fail to insert new proposal tags [%s]: %s", proposal.ProposalID, err)
	}

	err = p.cache.Upsert(proposal)

	if err != nil {
		log.Printf("fail to insert new proposal on cache [%s]: %s", proposal.ProposalID, err)
		return proposal.ProposalID, err
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
		return err
	}

	_, err = p.tags.Insert(proposal.Tags)

	if err != nil {
		log.Printf("fail to update proposal tags [%s]: %s", proposal.ProposalID, err)
	}

	err = p.cache.Upsert(proposal)

	if err != nil {
		log.Printf("fail to update new proposal on cache [%s]: %s", proposal.ProposalID, err)
		return err
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

	p.storage.InsertView(proposalID, "", "Load Proposal")

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
	if filter == nil {
		filter = &models.Filter{}
	}

	result, err := p.storage.Find(filter)
	//result, err := p.cache.Find(filter)

	if err != nil {
		log.Printf("fail to load data from filter: %s", err)
		return &models.ProposalsResponse{
			Filter: filter,
		}, err
	}

	//Não questione, o front pediu isso... :/
	tagMap := map[string]interface{}{}
	sideMap := map[models.Side]interface{}{}
	typeMap := map[models.Type]interface{}{}

	ids := make([]string, len(result))

	for i, p := range result {
		for _, t := range p.Tags {
			tagMap[t] = nil
		}
		sideMap[p.Side] = nil
		typeMap[p.ProposalType] = nil
		ids[i] = string(p.ProposalID)
	}

	p.storage.BulkInsertView(ids, "Load Proposal from Filter")

	tags := make([]string, len(tagMap))
	sides := make([]models.Side, len(sideMap))
	types := make([]models.Type, len(typeMap))

	i := 0
	for t := range tagMap {
		tags[i] = t
		i++
	}

	i = 0
	for s := range sideMap {
		sides[i] = s
		i++
	}

	i = 0
	for t := range typeMap {
		types[i] = t
		i++
	}
	//fim do pedido do front

	sort.Strings(tags)

	pgSize := int64(len(result))

	return &models.ProposalsResponse{
		Filter:              filter,
		Result:              result,
		ResultProposalTypes: types,
		ResultSides:         sides,
		ResultTags:          tags,
		CurrentPage:         &filter.PageNumber,
		CurrentPageSize:     &pgSize,
	}, err
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

	user, err := p.user.Load(string(prop.UserID))

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
					Contact:     formatPhones(user.Contact.Phones),
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
			case models.DataToShareTwitter:
				ret = append(ret, &models.DataToShareResponse{
					ContactType: models.DataToShareTwitter,
					Contact:     user.Contact.Twitter,
				})
			case models.DataToShareBankAccount:
				ret = append(ret, &models.DataToShareResponse{
					ContactType: models.DataToShareBankAccount,
					Contact:     formatBankAccs(prop),
				})
			}
		}
	}

	//TODO: Store user requested DTS on View
	p.storage.InsertView(proposalID, "", "DTS Request")

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

//InsertComplaint inserts a proposal complaint
func (p *Proposal) InsertComplaint(complaint *models.Complaint) error {
	if complaint != nil {
		log.Printf("A complaint happened: %v\n", *complaint)

		return p.storage.InsertComplaint(complaint)
	}

	return fmt.Errorf("cannot try insert an empty complaint")
}

func formatPhones(phones []*models.Phone) string {
	ret := make([]string, len(phones))

	for i, p := range phones {
		ret[i] = fmt.Sprintf("(%s) %s", p.Region, p.PhoneNumber)
	}

	return strings.Join(ret, ", ")
}

func formatBankAccs(proposal *models.Proposal) string {
	ret := []string{}

	if proposal == nil {
		return ""
	}

	if len(proposal.BankAccounts) == 0 {
		return ""
	}

	for _, acc := range proposal.BankAccounts {
		BranchDg := ""

		if len(acc.BranchDigit) > 0 {
			BranchDg = fmt.Sprintf("-%s", acc.BranchDigit)
		}

		ret = append(ret,
			fmt.Sprintf("Nome: %s\nDocumento: %s\nBco: (%d) %s\nAg: %s%s CC: %s-%s\n",
				acc.AccountOwner,
				acc.AccountDocument,
				acc.BankID,
				acc.BankName,
				acc.BranchNumber,
				BranchDg,
				acc.AccountNumber,
				acc.AccountDigit,
			),
		)
	}

	return strings.Join(ret, "\n")
}
