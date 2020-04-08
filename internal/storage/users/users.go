package users

//Users Object struct
type Users struct {
	conn *storage.Connection
}

//New creates a new instance
func New(conn *storage.Connection) *Users {
	return &Users{conn: conn}
}

const upsert = `
INSERT INTO 
USERS 
	(
		UserID,
		Name,
		Description,
		DeviceID,
		AllowShareData,
		Tags,
		Images,
		
		--Reputation
		Giver,
		Taker,
	
		--Contatct
		URL,
		Email,
		Facebook,
		Instagram,
		Google,
		AdditionalData,
		
		--Contact Address
		Address,
		City,
		State,
		ZipCode,
		Country,
	
		Lat,
		Long
	)
VALUES
   (
		$1, --UserID,
		$2, --Name,
		$3, --Description,
		$4, --DeviceID,
		$5, --AllowShareData,
		$6, --Tags,
		$7, --Images,
		 
		----Reputation
		$8, --Giver,
		$9, --Taker,
		 
		----Contatct
		$10, --URL,
		$11, --Email,
		$12, --Facebook,
		$13, --Instagram,
		$14, --Google,
		$15, --AdditionalData,
		 
		----Contact Address
		$16, --Address,
		$17, --City,
		$18, --State,
		$19, --ZipCode,
		$20, --Country,
		 
		$21, --Lat,
		$22 --Long
   ) 
ON CONFLICT (UserID) 
DO
	UPDATE
	SET 
		Name = $2,
		Description = $3,
		DeviceID = $4,
		AllowShareData = $5,
		Tags = $6,
		Images = $7,
		
		--Reputation
		Giver = $8,
		Taker = $9,

		--Contatct
		URL = $10,
		Email = $11,
		Facebook = $12,
		Instagram = $13,
		Google = $14,
		AdditionalData = $15,
		
		--Contact Address
		Address = $16,
		City = $17,
		State = $18,
		ZipCode = $19,
		Country = $20,

		Lat = $21,
		Long = $22;
`

//Upsert insert or update on database
func (u *users) Upsert(user *models.User) (string, error) {
	repGiver := 0.0
	repTaker := 0.0

	if user.Reputation != nil {
		retGiver = user.Reputation.Giver
		repTaker = user.Reputation.Taker
	}

	url := ""
	email := ""
	facebook := ""
	instagram := ""
	google := ""
	additionalData := map[string]interface{}

	if user.Contact != nil {
		url = user.Contact.URL
		email = user.Contact.Email
		facebook = user.Contact.Facebook
		instagram = user.Contact.Instagram
		google = user.Contact.Google				
		additionalData = user.Contact.AdditionalData
	}

	address := ""
	city := ""
	state := ""
	zipCode := 0
	country := ""

	lat := 0.0
	long := 0.0

	if user.Location != nil {
		address = user.Location.Address
		city := user.Location.City
		state := user.Location.State
		zipCode := user.Location.ZipCode
		country := user.Location.Country

		lat := user.Location.Lat
		long := user.Location.Long
	}

	aff, err := c.conn.Execute(
		upsert,
		user.UserID,
		user.Name,
		user.Description,
		user.AllowShareData,
		user.Tags,
		user.Images,
		repGiver,
		repTaker,
		url,
		email,
		facebook,
		instagram,
		google,
		additionalData,
		address,
		city,
		state,
		zipCode,
		country,
		lat,
		long,
	)

	if err == nil {
		err = u.insertPhones(user)
	}		

	return aff, err
}

const select = `
SELECT 
	UserID,
	Name,
	Description,
	DeviceID,
	AllowShareData,
	Tags,
	Images,

	--Reputation
	Giver,
	Taker,

	--Contatct
	URL,
	Email,
	Facebook,
	Instagram,
	Google,
	AdditionalData,

	--Contact Address
	Address,
	City,
	State,
	ZipCode,
	Country,

	Lat,
	Long
FROM
	USERS
WHERE
	UserID = $1;
`

//Load load from database
func (u *users) Load(userID string) (*models.User, error) {
	ret := models.User{
		UserID: userID,
		Contact: &models.Contact{},
		Reputation: &models.Reputation{}
		Location: &models.Location{}
	}

	row := c.conn.Db.QueryRow(select, userID)
	defer rows.Close()

	err := rows.Scan(
		&user.Name,
		&user.Description,
		&user.DeviceID,
		&user.AllowShareData,
		&user.Tags,
		&user.Images,
		&user.Reputation.Giver,
		&user.Reputation.Taker,
		&user.Contact.URL,
		&user.Contact.Email,
		&user.Contact.Facebook,
		&user.Contact.Instagram,
		&user.Contact.Google,
		&user.Contact.AdditionalData,
		&user.Location.Address,
		&user.Location.City,
		&user.Location.State,
		&user.Location.ZipCode,
		&user.Location.Country,
		&user.Location.Lat,
		&user.Location.Long,
	)

	ret.Contact.Phones = u.loadPhones(userID)

	return &ret, err
}

const insertPhones = `
INSERT INTO
PHONES
	UserID,
	CountryCode,
	IsDefault,
	PhoneNumber,
	Region,
	WhastApp
VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
	);
`

const cleanPhones = `
DELETE FROM 
	PHONES
WHERE
	UserID = $1
`

func (u *Users) insertPhones(user *models.Users) error {
	if user!= nil {
		if user.Contact != nil {
			if len(user.Contact.Phones) > 0 {
				_, err := c.conn.Db.Execute(cleanPhones, user.UserID)

				if err != nil {
					return err
				}

				for _, p := range user.Contact.Phones {
					qtd, err := c.conn.Db.Execute(
						insertPhones, 
						user.UserID,
						p.CountryCode,
						p.IsDefault,
						p.PhoneNumber,
						p.Region,
						p.WhastApp,
					)

					if err != nil {
						return err
					}
				}
			}
		}		
	}

	return err
}

const selectPhones = `
SELECT 
	CountryCode,
	IsDefault,
	PhoneNumber,
	Region,
	WhastApp
FROM
	PHONES
WHERE
	UserID = $1
ORDER BY
	IsDefault DESC,
	CreatedAt
;
`

func (u *Users) loadPhones(userID string) ([]*models.Phone, error) {
	ret := []*models.Phone{}

	rows, err := c.conn.Db.QueryRow(select, userID)

	if err != nil {
		return ret, err
	}

	defer rows.Close()

	for rows.Next() {
		p &models.Phone{}

		err = rows.Scan(
			&p.CountryCode,
			&p.IsDefault,
			&p.PhoneNumber,
			&p.Region,
			&p.WhastApp,
		)

		if err != nil {
			return ret, err
		}

		ret = append(ret, p)
	}

	return ret, err
}