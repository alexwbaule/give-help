package etc

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
)

//Tags Object struct
type Etc struct {
	conn *connection.Connection
}

//New creates a new instance
func New(conn *connection.Connection) *Etc {
	return &Etc{conn: conn}
}

const selectEtc = `
SELECT 
	LOWER(Key) as Key,
	Value
FROM 
	ETC
WHERE
	IsActive = true
ORDER BY Key;`

//Load load etc key value from database
func (e *Etc) Load() (models.Etc, error) {
	ret := make(map[string]string)

	db := e.conn.Get()

	rows, err := db.Query(selectEtc)

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			k := ""
			v := ""
			if err = rows.Scan(&k, &v); err == nil {
				ret[k] = v
			}
		}
	}

	return models.Etc(ret), e.conn.CheckError(err)
}
