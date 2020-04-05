package common

type User struct {
	UserId         string            `json:",omitempty"`
	CreatedAt      time.time         `json:",omitempty"`
	Name           string            `json:",omitempty"`
	Description    string            `json:",omitempty"`
	DeviceId       string            `json:",omitempty"`
	Location       *Location         `json:",omitempty"`
	Contact        *Contact          `json:",omitempty"`
	AllowShareData bool              `json:",omitempty"`
	Images         []string          `json:",omitempty"`
	Reputation     *Reputation       `json:",omitempty"`
	Tags           map[string]string `json:",omitempty"`
}

type Contact struct {
	Phones         []*Phone
	Email          string
	Instagram      string            `json:",omitempty"`
	Facebook       string            `json:",omitempty"`
	Google         string            `json:",omitempty"`
	Url            string            `json:",omitempty"`
	AddtiionalData map[string]string `json:",omitempty"`
}

type Phone struct {
	Region      string `json:",omitempty"`
	Number      string `json:",omitempty"`
	CountryCode string `json:",omitempty"`
	WhatsApp    bool   `json:",omitempty"`
	Default     bool   `json:",omitempty"`
}

type Location struct {
	ZipCode int     `json:",omitempty"`
	Address string  `json:",omitempty"`
	City    string  `json:",omitempty"`
	State   string  `json:",omitempty"`
	Country string  `json:",omitempty"`
	Lat     float64 `json:",omitempty"`
	Long    float64 `json:",omitempty"`
}

type Reputation struct {
	Giver float64 `json:",omitempty"`
	Taker float64 `json:",omitempty"`
}
