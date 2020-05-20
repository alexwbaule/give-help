package proposal

import (
	"context"
	"fmt"
	"strings"

	"log"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	"github.com/lib/pq"
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
	Title,
    Description,
	ProposalValidate,
	City,
	State,
	Country,	
    Lat,
    Long,
	Range,
	AreaTags,
	IsActive,	
	Images,
	DataToShare,
	ExposeUserData,
	EstimatedValue,
	Ranking
) 
VALUES
(
	$1, --ProposalID,
    $2, --UserID,
    $3, --Side,
    $4, --ProposalType,
	$5, --Tags
	$6, --Title
    $7, --Description
	$8, --ProposalValidate,
	$9, --City,
	$10, --State,
	$11, --Country,
    $12, --Lat,
    $13, --Long,
    $14, --Range,
	$15, --AreaTags,
	$16, --IsActive,
	$17, --Images,
	$18, --DataToShare,
	$19, --ExposeUserData,
	$20, --EstimatedValue,
	$21  --Ranking
)
ON CONFLICT (ProposalID) 
DO UPDATE SET
	LastUpdate = CURRENT_TIMESTAMP,
    Side = $3,
    ProposalType = $4,
	Tags = $5,
	Title = $6,
    Description = $7,
	ProposalValidate = $8,
	City = $9,
	State = $10,
	Country = $11,	
    Lat = $12,
    Long = $13,
	Range = $14,
	AreaTags = $15,
	IsActive = $16,	
	Images = $17,
	DataToShare = $18,
	ExposeUserData = $19,
	EstimatedValue = $20,
	Ranking = $21
;
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

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	lat := float64(0)
	long := float64(0)
	areaRange := float64(0)
	areaTags := []string{}
	city := ""
	state := ""
	country := ""

	if proposal.TargetArea != nil {
		city = proposal.TargetArea.City
		state = proposal.TargetArea.State
		country = proposal.TargetArea.Country

		lat = *proposal.TargetArea.Lat
		long = *proposal.TargetArea.Long
		areaRange = proposal.TargetArea.Range

		areaTags = common.NormalizeTagArray(proposal.TargetArea.AreaTags)
	}

	if lat == 0 {
		lat = -23.5486
	}

	if long == 0 {
		long = -46.6392
	}

	if proposal.DataToShare == nil {
		proposal.DataToShare = []models.DataToShare{}
	}

	if proposal.Images == nil {
		proposal.Images = []string{}
	}

	_, err = db.ExecContext(
		ctx,
		upsertProposal,
		proposal.ProposalID,
		proposal.UserID,
		proposal.Side,
		proposal.ProposalType,
		pq.Array(common.NormalizeTagArray(proposal.Tags)),
		proposal.Title,
		proposal.Description,
		proposal.ProposalValidate,
		city,
		state,
		country,
		lat,
		long,
		areaRange,
		pq.Array(common.NormalizeTagArray(areaTags)),
		proposal.IsActive,
		pq.Array(proposal.Images),
		pq.Array(proposal.DataToShare),
		proposal.ExposeUserData,
		proposal.EstimatedValue,
		proposal.Ranking,
	)

	if err != nil {
		if perr, ok := err.(*pq.Error); ok {
			tx.Rollback()
			return fmt.Errorf("fail to try execute upsert proposal data: proposal=%v pq-error=%s", proposal, perr)
		}

		tx.Rollback()
		return fmt.Errorf("fail to try execute upsert proposal data: proposal=%v error=%s", proposal, err)
	}

	err = p.upsertAccounts(ctx, string(proposal.ProposalID), proposal.BankAccounts)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("fail to try execute upsert proposal bank accounts: proposal=%v error=%s", proposal, err)
	}

	return tx.Commit()
}

const insertAccounts = `INSERT INTO BANK_ACCOUNTS 
(
	AccountNumber,
	AccountDigit,
	AccountOwner,
	AcountDocument,
	BranchNumber,
	BranchDigit,
	BankID,
	ProposalID
) 
VALUES 
(
	$1, --AccountNumber,
	$2, --AccountDigit,
	$3, --AccountOwner,
	$4, --AcountDocument,
	$5, --BranchNumber,
	$6, --BranchDigit
	$7, --BankID,
	$8 --ProposalID
);
`

const removeProposalAccounts = `
DELETE FROM BANK_ACCOUNTS WHERE ProposalID = $1;
`

func (p *Proposal) upsertAccounts(ctx context.Context, proposalID string, accs []*models.BankAccount) error {
	db := p.conn.Get()

	//clean proposal accounts
	if _, err := db.ExecContext(ctx, removeProposalAccounts, proposalID); err != nil {
		log.Printf("fail to try clean proposal bank accounts, calling rollback: %s", err)
		return p.conn.CheckError(err)
	}

	for _, acc := range accs {
		result, err := db.ExecContext(
			ctx,
			insertAccounts,
			acc.AccountNumber,
			acc.AccountDigit,
			acc.AccountOwner,
			acc.AccountDocument,
			acc.BranchNumber,
			acc.BranchDigit,
			acc.BankID,
			proposalID,
		)

		if err != nil {
			log.Printf("fail to try insert new proposal bank accounts (insert fail), calling rollback: %s", err)
			return p.conn.CheckError(err)
		}

		aff, err := result.RowsAffected()

		if err != nil {
			log.Printf("fail to try insert new proposal bank accounts (read result), calling rollback: %s", err)
			return p.conn.CheckError(err)
		}

		if aff == 0 {
			log.Printf("fail to try insert new proposal bank accounts (no rows affected), calling rollback: %s", err)
			return fmt.Errorf("0 rows affected, check arguments!")
		}
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
	Title,
	Description,
	ProposalValidate,
	City,
	State,
	Country,
	Lat,
	Long,
	Range,
	AreaTags,
	IsActive,
	Images,
	EstimatedValue,
	ExposeUserData,
	DataToShare,
	Ranking
FROM
	PROPOSALS
WHERE	
	%s
ORDER BY
	Ranking DESC,
	CreatedAt ASC
%s
`

//LoadFromProposal load an unique proposal from a proposalID
func (p *Proposal) LoadFromID(proposalID string) (*models.Proposal, error) {
	ret := models.Proposal{
		TargetArea: &models.Location{
			City:    "",
			State:   "",
			Country: "",
		},
		DataToShare: []models.DataToShare{},
	}

	cmd := fmt.Sprintf(selectProposal, "ProposalID = $1", "")

	db := p.conn.Get()

	var tags []string
	var areaTags []string
	var images []string
	var dataToShare []string
	var city *string
	var state *string
	var country *string

	err := db.QueryRow(cmd, proposalID).Scan(
		&ret.ProposalID,
		&ret.UserID,
		&ret.CreatedAt,
		&ret.LastUpdate,
		&ret.Side,
		&ret.ProposalType,
		pq.Array(&tags),
		&ret.Title,
		&ret.Description,
		&ret.ProposalValidate,
		&city,
		&state,
		&country,
		&ret.TargetArea.Lat,
		&ret.TargetArea.Long,
		&ret.TargetArea.Range,
		pq.Array(&areaTags),
		&ret.IsActive,
		pq.Array(&images),
		&ret.EstimatedValue,
		&ret.ExposeUserData,
		pq.Array(&dataToShare),
		&ret.Ranking,
	)

	if city != nil {
		ret.TargetArea.City = string(*city)
	}

	if state != nil {
		ret.TargetArea.State = string(*state)
	}

	if country != nil {
		ret.TargetArea.Country = string(*country)
	}

	ret.Tags = common.NormalizeTagArray(tags)
	ret.TargetArea.AreaTags = common.NormalizeTagArray(areaTags)
	ret.Images = images

	shareBankAcc := false

	ret.DataToShare = make([]models.DataToShare, len(dataToShare))
	for i, v := range dataToShare {
		ret.DataToShare[i] = models.DataToShare(v)

		if ret.DataToShare[i] == models.DataToShareBankAccount {
			shareBankAcc = true
		}
	}

	if shareBankAcc {
		ret.BankAccounts, err = p.loadAccounts(proposalID)

		if err != nil {
			log.Printf("fail to try load proposal bank accounts - error: %s", err)
			return &ret, p.conn.CheckError(err)
		}
	}

	return &ret, p.conn.CheckError(err)
}

//LoadFromUser load all proposals from an userID
func (p *Proposal) LoadFromUser(userID string) ([]*models.Proposal, error) {
	cmd := fmt.Sprintf(selectProposal, "UserID = $1", "")

	return p.load(cmd, userID)
}

//Find find all proposals that match with filter
func (p *Proposal) Find(filter *models.Filter) ([]*models.Proposal, error) {
	if filter == nil {
		filter = &models.Filter{
			PageSize:        50,
			PageNumber:      0,
			IncludeExpired:  false,
			IncludeInactive: false,
		}
	}

	args := []interface{}{}
	wheres := []string{}

	if len(filter.Description) > 0 {
		likeTarget := "%" + strings.ToLower(strings.TrimSpace(filter.Description)) + "%"
		args = append(args, likeTarget)
		wheres = append(wheres, fmt.Sprintf("( LOWER(CONCAT(Description, Title, array_to_string(AreaTags, ','), array_to_string(Tags, ','), City, State, Country)) LIKE $%d ) ", len(args)))
	}

	if len(filter.UserID) > 0 {
		args = append(args, filter.UserID)
		wheres = append(wheres, fmt.Sprintf("UserID = $%d", len(args)))
	}

	if len(filter.Side) > 0 {
		args = append(args, filter.Side)
		wheres = append(wheres, fmt.Sprintf("Side = $%d", len(args)))
	}

	for _, t := range filter.ProposalTypes {
		args = append(args, t)
		wheres = append(wheres, fmt.Sprintf("ProposalType = $%d", len(args)))
	}

	for _, t := range filter.Tags {
		args = append(args, "%"+strings.ToLower(strings.TrimSpace(t))+"%")
		wheres = append(wheres, fmt.Sprintf("array_to_string(Tags, ',') LIKE $%d", len(args)))
	}

	if filter.MaxValue != nil && *filter.MaxValue > 0 {
		args = append(args, filter.MinValue, filter.MaxValue)
		wheres = append(wheres, fmt.Sprintf("(EstimatedValue >= $%d AND EstimatedValue <= $%d )", len(args)-1, len(args)))
	}

	if filter.MinValue != nil && *filter.MinValue > 0 {
		args = append(args, filter.MinValue)
		wheres = append(wheres, fmt.Sprintf("EstimatedValue >= $%d", len(args)))
	}

	if filter.TargetArea != nil {
		for _, t := range filter.TargetArea.AreaTags {
			args = append(args, "%"+strings.ToLower(strings.TrimSpace(t))+"%")
			wheres = append(wheres, fmt.Sprintf("array_to_string(AreaTags, ',') LIKE $%d", len(args)))
		}

		rang := filter.TargetArea.Range

		if filter.TargetArea.Lat != nil && filter.TargetArea.Long != nil {
			if rang < 1 {
				rang = 1
			}

			if N, S, W, E, err := common.CalculeRange(filter.TargetArea); err == nil {
				args = append(args, S, N, E, W)
				wheres = append(
					wheres,
					fmt.Sprintf(
						`( (Lat BETWEEN $%d AND $%d) AND (Long BETWEEN $%d AND $%d ) )`,
						len(args)-3,
						len(args)-2,
						len(args)-1,
						len(args)),
				)
			}
		} else {
			log.Printf("cannot calculate target area range filter")
		}
	}

	andFilters := []string{} // "ProposalValidate >= %s ", "IsActive = true" }

	if !filter.IncludeExpired {
		andFilters = append(andFilters, "ProposalValidate >= CURRENT_TIMESTAMP")
	}

	if !filter.IncludeInactive {
		andFilters = append(andFilters, "IsActive = true")
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 50
	}

	if filter.PageNumber < 0 {
		filter.PageNumber = 0
	}

	var baseFilter string

	switch len(andFilters) {
	case 0:
		baseFilter = ""
	case 1:
		baseFilter = andFilters[0]
	default:
		baseFilter = strings.Join(andFilters, " AND ")
	}

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", filter.PageSize, filter.PageNumber*filter.PageSize)
	var cmd string
	if len(wheres) > 0 {
		cmd = fmt.Sprintf(
			selectProposal,
			fmt.Sprintf(
				"( %s ) \n\tAND ( %s )",
				strings.Join(wheres, " OR \n\t"),
				baseFilter,
			),
			limit,
		)
	} else {
		cmd = fmt.Sprintf(
			selectProposal,
			baseFilter,
			limit,
		)
	}

	return p.load(cmd, args...)
}

func (p *Proposal) load(cmd string, args ...interface{}) ([]*models.Proposal, error) {
	ret := []*models.Proposal{}

	db := p.conn.Get()

	rows, err := db.Query(cmd, args...)

	if err != nil {
		return ret, p.conn.CheckError(err)
	}

	defer rows.Close()

	for rows.Next() {
		i := models.Proposal{
			TargetArea: &models.Location{
				City:    "",
				State:   "",
				Country: "",
			},
		}

		var tags []string
		var areaTags []string
		var images []string
		var dataToShare []string
		var city *string
		var state *string
		var country *string

		err = rows.Scan(
			&i.ProposalID,
			&i.UserID,
			&i.CreatedAt,
			&i.LastUpdate,
			&i.Side,
			&i.ProposalType,
			pq.Array(&tags),
			&i.Title,
			&i.Description,
			&i.ProposalValidate,
			&city,
			&state,
			&country,
			&i.TargetArea.Lat,
			&i.TargetArea.Long,
			&i.TargetArea.Range,
			pq.Array(&areaTags),
			&i.IsActive,
			pq.Array(&images),
			&i.EstimatedValue,
			&i.ExposeUserData,
			pq.Array(&dataToShare),
			&i.Ranking,
		)

		if city != nil {
			i.TargetArea.City = string(*city)
		}

		if state != nil {
			i.TargetArea.State = string(*state)
		}

		if country != nil {
			i.TargetArea.Country = string(*country)
		}

		if err != nil {
			fmtArgs := make([]string, len(args))
			for p, a := range args {
				fmtArgs[p] = fmt.Sprintf("$%d=%v", p, a)
			}
			log.Printf("query error. \nquery: %s \nargs: %s\nerror: %s", cmd, strings.Join(fmtArgs, ";"), err)
			return ret, p.conn.CheckError(err)
		}

		i.Tags = common.NormalizeTagArray(tags)
		i.TargetArea.AreaTags = common.NormalizeTagArray(areaTags)
		i.Images = images

		shareBankAcc := false

		i.DataToShare = make([]models.DataToShare, len(dataToShare))
		for pos, v := range dataToShare {
			i.DataToShare[pos] = models.DataToShare(v)

			if i.DataToShare[pos] == models.DataToShareBankAccount {
				shareBankAcc = true
			}
		}

		if shareBankAcc {
			i.BankAccounts, err = p.loadAccounts(string(i.ProposalID))

			if err != nil {
				log.Printf("fail to try load proposal bank accounts - error: %s", err)
				return ret, p.conn.CheckError(err)
			}
		}

		ret = append(ret, &i)
	}

	if err != nil {
		fmtArgs := make([]string, len(args))
		for p, a := range args {
			fmtArgs[p] = fmt.Sprintf("$%d=%v", p, a)
		}
		log.Printf("query error. \nquery: %s \nargs: %s\nerror: %s", cmd, strings.Join(fmtArgs, ";"), err)
	}

	return ret, p.conn.CheckError(err)
}

const selectAccounts = `
SELECT 
	B.BankID,
	B.BankName,
	B.BankFullName,
	A.AccountNumber,
	A.AccountDigit,
	A.AccountOwner,
	A.AcountDocument,
	A.BranchNumber,
	A.BranchDigit
FROM 
	BANK_ACCOUNTS A INNER JOIN BANKS B
		ON A.BankID = B.BankID
WHERE
	A.ProposalID = $1
ORDER BY 
	A.CreatedAt;
`

func (p *Proposal) loadAccounts(proposalId string) ([]*models.BankAccount, error) {
	ret := []*models.BankAccount{}

	db := p.conn.Get()

	rows, err := db.Query(selectAccounts, proposalId)

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			acc := &models.BankAccount{}

			if err = rows.Scan(
				&acc.BankID,
				&acc.BankName,
				&acc.BankFullname,
				&acc.AccountNumber,
				&acc.AccountDigit,
				&acc.AccountOwner,
				&acc.AccountDocument,
				&acc.BranchNumber,
				&acc.BranchDigit,
			); err == nil {
				ret = append(ret, acc)
			} else {
				return ret, p.conn.CheckError(err)
			}
		}
	}

	return ret, p.conn.CheckError(err)
}

const insertComplaint = `INSERT INTO COMPLAINTS
(
	Complainer,
	ProposalID,
	Comment,
	Accepted
) 
VALUES 
(
	$1,
	$2,
	$3,
	false
);`

//InsertComplaint insert categories on database
func (p *Proposal) InsertComplaint(complaint *models.Complaint) error {
	db := p.conn.Get()

	if len(complaint.ProposalID) == 0 {
		return fmt.Errorf("fail to try insert a complaint: not allowed empty proposalID: %v", *complaint)
	}

	if len(complaint.Comment) == 0 {
		return fmt.Errorf("fail to try insert a complaint: not allowed empty comment: %v", *complaint)
	}

	if len(complaint.Complainer) == 0 {
		complaint.Complainer = "system:anonymous"
	}

	result, err := db.Exec(insertComplaint, complaint.Complainer, complaint.ProposalID, complaint.Comment)

	if err != nil {
		log.Printf("fail to try insert a complaint: %v", *complaint)
		return fmt.Errorf("fail to try insert a complaint: %v", *complaint)
	}

	if aff, _ := result.RowsAffected(); aff == 0 {
		log.Printf("fail to try insert a complaint: %v", *complaint)
		return fmt.Errorf("fail to try insert a complaint: %v", *complaint)
	}

	return p.conn.CheckError(err)
}
