package storage

func (s *Storage) InsertProposal(proposal *models.Proposal) error {
	payload, err := json.Marshal(proposal)

	if err != nil {
		return err
	}

	_, err = s.send(payload)

	return err
}

func (s *Storage) LoadProposal(userID string) (*models.Proposal, error) {

}

func (s *Storage) FindProposals(lat float64, long float64, range float64, tags []string) ([]*models.Proposal, error) {

}

func (s *Storega) ChangeProposalActiveStatus(status bool) error {
	
}

func (s *Storega) ChangeProposalValidate(validate time.Time) error {
	
}
