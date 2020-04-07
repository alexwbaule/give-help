package categories

type Categories struct {
	conn *storage.Connection
}

func New(conn *storage.Connection) *Categories {
	return &Categories{conn: conn}
}

func (c *Categories) Insert(categories []string) error {
	return nil
}

func (c *Categories) Load() ([]string, error) {
	return []string{""}, nil
}
