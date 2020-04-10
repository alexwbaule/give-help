package proposal

import (
	"fmt"
	"strings"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	"github.com/lib/pq"
	"log"
)

//Proposal Object struct
type Proposal struct {
	conn *connection.Connection
}

//New creates a new instance
func New(conn *connection.Connection) *Proposal {
	return &Proposal{conn: conn}
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
	LastUpdate = CURRENT_TIMESTAMP,
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

//Upsert insert or update on database
func (p *Proposal) Upsert(proposal *models.Proposal) error {
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

		for _, t := range proposal.TargetArea.AreaTags {
			areaTags = append(areaTags, strings.ToUpper(t))
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

//LoadFromProposal load an unique proposal from a proposalID
func (p *Proposal) LoadFromProposal(prposalID string) (*models.Proposal, error) {
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

//LoadFromUser load all proposals from an userID
func (p *Proposal) LoadFromUser(userID string) ([]*models.Proposal, error) {
	cmd := fmt.Sprintf(selectProposal, "UserID = $1")

	return p.load(cmd, userID)
}

//Find find all proposals that match with filter
func (p *Proposal) Find(filter *models.Filter) ([]*models.Proposal, error) {
	if filter == nil {
		log.Printf("cannot execute a query with null filter")
		return nil, nil
	}

	args := []interface{}{}
	wheres := []string{}

	if len(filter.Description) > 0 {
		args = append(args, "%"+strings.ToUpper(filter.Description)+"%")
		wheres = append(wheres, fmt.Sprintf("UPPER(Description) LIKE $%d", len(wheres)+1))
	}

	if len(filter.Side) > 0 {
		args = append(args, filter.Side)
		wheres = append(wheres, fmt.Sprintf("Side = $%d", len(wheres)+1))
	}

	for _, t := range filter.ProposalTypes {
		args = append(args, t)
		wheres = append(wheres, fmt.Sprintf("ProposalType = $%d", len(wheres)+1))
	}

	for _, t := range filter.Tags {
		args = append(args, strings.ToUpper(t))
		wheres = append(wheres, fmt.Sprintf("$%d = ANY(Tags)", len(wheres)+1))
	}

	if filter.TargetArea != nil {
		for _, t := range filter.TargetArea.AreaTags {
			args = append(args, strings.ToUpper(t))
			wheres = append(wheres, fmt.Sprintf("$%d = ANY(AreaTags)", len(wheres)+1))
		}

		rang := filter.TargetArea.Range

		if filter.TargetArea.Lat != 0 && filter.TargetArea.Long != 0 {
			if rang < 1 {
				rang = 1
			}

			if N, S, W, E, err := common.CalculeRange(filter.TargetArea); err == nil {
				args = append(args, S)
				args = append(args, N)
				args = append(args, E)
				args = append(args, W)
				wheres = append(
					wheres,
					fmt.Sprintf(
						`( (Lat BETWEEN $%d AND $%d) AND (Long BETWEEN $%d AND $%d ) )`,
						len(wheres)+1,
						len(wheres)+2,
						len(wheres)+3,
						len(wheres)+4),
				)
			}
		} else {
			log.Printf("cannot calculate target area range filter")
		}
	}

	cmd := fmt.Sprintf(selectProposal, strings.Join(wheres, " OR \n\t")+"\n")

	return p.load(cmd, args...)
}

func (p *Proposal) load(cmd string, args ...interface{}) ([]*models.Proposal, error) {
	ret := []*models.Proposal{}

	db := p.conn.Get()
	defer db.Close()

	rows, err := db.Query(cmd, args...)

	if err != nil {
		return ret, p.conn.CheckError(err)
	}

	defer rows.Close()

	for rows.Next() {
		i := models.Proposal{TargetArea: &models.Area{}}
		var tags []string
		var areaTags []string

		err = rows.Scan(
			&i.ProposalID,
			&i.UserID,
			&i.CreatedAt,
			&i.LastUpdate,
			&i.Side,
			&i.ProposalType,
			pq.Array(&tags),
			&i.Description,
			&i.ProposalValidate,
			&i.TargetArea.Lat,
			&i.TargetArea.Long,
			&i.TargetArea.Range,
			pq.Array(&areaTags),
			&i.IsActive,
		)

		if err != nil {
			return ret, p.conn.CheckError(err)
		}

		i.Tags = tags
		i.TargetArea.AreaTags = areaTags

		ret = append(ret, &i)
	}

	return ret, p.conn.CheckError(err)
}
