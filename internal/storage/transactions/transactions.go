package transactions

type Transactions struct {
	conn *storage.Connection
}

func New(conn *storage.Connection) *Transactions {
	return &Transactions{conn: conn}
}

func (t *Transactions) Insert(transaction *models.Transaction) error {

}

func (t *Transactions) InsertGiverReview(transactionID string, review *models.Review) error {

}

func (t *Transactions) InsertTakerReview(transactionID string, review *models.Review) error {

}

func (t *Transactions) UpdatStatus(transactionID string, status models.TransactionStatus) error {

}
