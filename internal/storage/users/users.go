package users

type Users struct {
	conn *storage.Connection
}

func New(conn *storage.Connection) *Users {
	return &Users{conn: conn}
}

func (u *users) Upsert(user *models.User) error {

}

func (u *users) Load(userID string) (*models.User, error) {

}
