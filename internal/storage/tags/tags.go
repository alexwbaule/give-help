package tags

import (
	"fmt"
	"strings"

	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

//Tags Object struct
type Tags struct {
	conn *connection.Connection
}

//New creates a new instance
func New(conn *connection.Connection) *Tags {
	return &Tags{conn: conn}
}

const insertCategories = `INSERT INTO TAGS 
(
	Name
) 
VALUES 
%s
ON CONFLICT (Name) 
DO NOTHING;`

//Insert insert categories on database
func (t *Tags) Insert(categories []string) (int64, error) {
	items := make([]string, len(categories))
	for pos, cat := range categories {
		if len(cat) > 0 {
			items[pos] = fmt.Sprintf(`('%s')`, cat)
		}
	}

	if len(items) == 0 {
		return 0, nil
	}

	cmd := fmt.Sprintf(insertCategories, strings.Join(items, ","))

	db := t.conn.Get()

	aff, err := db.Exec(cmd)

	if err != nil {
		return 0, t.conn.CheckError(err)
	}

	return aff.RowsAffected()
}

const selectCategories = `
SELECT 
	DISTINCT Name
FROM 
	TAGS
WHERE
	Name IS NOT NULL
	AND LENGTH(Name) > 0
ORDER BY NAME`

//Load load categories from database
func (t *Tags) Load() ([]string, error) {
	ret := []string{}

	db := t.conn.Get()

	rows, err := db.Query(selectCategories)

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			var name string
			if err = rows.Scan(&name); err == nil {
				ret = append(ret, name)
			}
		}
	}

	return ret, t.conn.CheckError(err)
}
