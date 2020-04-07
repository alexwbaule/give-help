package proposals

type Proposals struct {
	conn *storage.Connection
}

func New(conn *storage.Connection) *Proposals {
	return &Proposals{conn: conn}
}

func (p *Proposals) Insert(proposal *models.Proposal) error {

}

func (p *Proposals) Load(userID string) (*models.Proposal, error) {

}

func (p *Proposals) Find(lat float64, long float64, range float64, tags []string) ([]*models.Proposal, error) {

}

func (p *Proposals) ChangeActiveStatus(status bool) error {
	
}

func (p *Proposals) ChangeValidate(validate time.Time) error {
	
}
