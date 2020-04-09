package proposals

import (
	"fmt"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage"
	"github.com/lib/pq"
)

type Proposals struct {
	conn *storage.Connection
}

func New(conn *storage.Connection) *Proposals {
	return &Proposals{conn: conn}
}

const upsertProposal string = `
INSERT INTO PROPOSALS (
	ProposalID,
    UserID,
    Side,
    ProposalType,
    Tags,
    Description,
    ProposalValidate,
    Lat,
    Long,
	Range,
	AreaTags,
    IsActive
) 
VALUES
(
	$1, --ProposalID,
    $2, --UserID,
    $3, --Side,
    $4, --ProposalType,
    $5, --Tags
    $6, --Description
    $7, --ProposalValidate,
    $8, --Lat,
    $9, --Long,
    $10, --Range,
	$11, --AreaTags
    $12 --IsActive
)
ON CONFLICT (ProposalID) 
DO UPDATE SET
	LastUpdate = CURRENT_DATE,
    Side = $3,
    ProposalType = $4,
    Tags = $5,
    Description = $6,
    ProposalValidate = $7,
    Lat = $8,
    Long = $9,
	Range = $10,
	AreaTags = $11,
    IsActive = $12;
`

func (p *Proposals) Upsert(proposal *models.Proposal) error {
	if proposal == nil {
		return fmt.Errorf("cannot insert an empty proposal struct")
	}

	if len(proposal.UserID) == 0 {
		return fmt.Errorf("cannot insert an empty UserID")
	}

	if len(proposal.ProposalID) == 0 {
		return fmt.Errorf("cannot insert an empty ProposalID")
	}

	db := p.conn.Get()
	defer db.Close()

	lat := float64(0)
	long := float64(0)
	areaRange := float64(0)
	areaTags := []string{}

	if proposal.TargetArea != nil {
		lat = proposal.TargetArea.Lat
		long = proposal.TargetArea.Long
		areaRange = proposal.TargetArea.Range

		if len(proposal.TargetArea.AreaTags) > 0 {
			areaTags = proposal.TargetArea.AreaTags
		}
	}

	_, err := db.Exec(
		upsertProposal,
		proposal.ProposalID,
		proposal.UserID,
		proposal.Side,
		proposal.ProposalType,
		pq.Array(proposal.Tags),
		proposal.Description,
		proposal.ProposalValidate,
		lat,
		long,
		areaRange,
		pq.Array(areaTags),
		proposal.IsActive,
	)

	if err != nil {
		if perr, ok := err.(*pq.Error); ok {
			return fmt.Errorf("fail to try execute upsert proposal data: proposal=%v pq-error=%s", proposal, perr)
		}

		return fmt.Errorf("fail to try execute upsert proposal data: proposal=%v error=%s", proposal, err)
	}

	return nil
}

const selectProposal string = `
SELECT
	ProposalID,
	UserID,
	CreatedAt,
	LastUpdate,
	Side,
	ProposalType,
	Tags,
	Description,
	ProposalValidate,
	Lat,
	Long,
	Range,
	AreaTags,
	IsActive
FROM
	PROPOSALS
WHERE	
	%s
ORDER BY
	CreatedAt ASC

`

func (p *Proposals) LoadFromProposal(prposalID string) (*models.Proposal, error) {
	ret := models.Proposal{TargetArea: &models.Area{}}

	cmd := fmt.Sprintf(selectProposal, "ProposalID = $1")

	db := p.conn.Get()
	defer db.Close()

	var tags []string
	var areaTags []string

	err := db.QueryRow(cmd, prposalID).Scan(
		&ret.ProposalID,
		&ret.UserID,
		&ret.CreatedAt,
		&ret.LastUpdate,
		&ret.Side,
		&ret.ProposalType,
		pq.Array(&tags),
		&ret.Description,
		&ret.ProposalValidate,
		&ret.TargetArea.Lat,
		&ret.TargetArea.Long,
		&ret.TargetArea.Range,
		pq.Array(&areaTags),
		&ret.IsActive,
	)

	ret.Tags = tags
	ret.TargetArea.AreaTags = areaTags

	return &ret, p.conn.CheckError(err)
}

func (p *Proposals) LoadFromUser(userID string) ([]*models.Proposal, error) {
	ret := []*models.Proposal{}

	cmd := fmt.Sprintf(selectProposal, "UserID = $1")

	db := p.conn.Get()
	defer db.Close()

	rows, err := db.Query(cmd, userID)

	if err != nil {
		return ret, p.conn.CheckError(err)
	}

	defer rows.Close()

	for rows.Next() {
		i := models.Proposal{TargetArea: &models.Area{}}

		err = rows.Scan(
			&i.ProposalID,
			&i.UserID,
			&i.CreatedAt,
			&i.LastUpdate,
			&i.Side,
			&i.ProposalType,
			pq.Array(&i.Tags),
			&i.Description,
			&i.ProposalValidate,
			&i.TargetArea.Lat,
			&i.TargetArea.Long,
			&i.TargetArea.Range,
			pq.Array(&i.TargetArea.AreaTags),
			&i.IsActive,
		)

		if err != nil {
			return ret, p.conn.CheckError(err)
		}

		ret = append(ret, &i)
	}

	return ret, p.conn.CheckError(err)
}

func (p *Proposals) Find(filter *models.Filter) ([]*models.Proposal, error) {
	return nil, nil
}

func (p *Proposals) ChangeActiveStatus(status bool) error {
	return nil
}

func (p *Proposals) ChangeValidate(validate time.Time) error {
	return nil
}
