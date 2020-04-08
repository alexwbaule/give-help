package categories

//Categories Object struct
type Categories struct {
	conn *storage.Connection
}

//New creates a new instance
func New(conn *storage.Connection) *Categories {
	return &Categories{conn: conn}
}

const insert = `INSERT INTO CATEGORIES (Name) VALUES %s;`

//Insert insert categories on database
func (c *Categories) Insert(categories []string) (int, error) {
	items := []string{}
	for pos, cat := range categories {
		if len(cat) > 0 {
			items = append(items, fmt.Sprintf(`('%s')`, cat))
		}		
	}

	if len(items) == 0 {
		return 0, nil
	}

	cmd := fmt.Sprintf(insert, strings.Join(items, ","))

	return c.conn.Execute(cmd)
}

const select = `
SELECT 
	DISTINCT Name 
FROM 
	CATEGORIES
WHERE
	Name IS NOT NULL
	AND (LENGTH) > 0
ORDER BY NAME`

//Load load categories from database
func (c *Categories) Load() ([]string, error) {
	ret := []string{}

	rows, err := c.conn.Db.Query(select)
	
	if err == nil {
		for rows.Next() {
			var name string
			if err = rows.Scan(&name); err == nil {
				ret = append(ret, name)
			}
		}
	}

	return ret, err
}
