package common

import "time"

type Side string

const (
	Offer         Side = "offer"
	Request            = "request"
	LocalBusiness      = "local-business"
)

type Type string

const (
	Job     Type = "job"
	Service      = "service"
	Product      = "product"
	Finance      = "finance"
)

type Transaction struct {
	TransactionId string            `json:",omitempty"`
	CreatedAt     time.time         `json:",omitempty"`
	Side          Side              `json:",omitempty"`
	Type          Type              `json:",omitempty"`
	UserId        string            `json:",omitempty"`
	Tags          map[string]string `json:",omitempty"`
	Description   string            `json:",omitempty"`
	Validate      time.Time         `json:",omitempty"`
	TargetArea    *Area             `json:",omitempty"`
}

type Area struct {
	Lat   float64           `json:",omitempty"`
	Long  float64           `json:",omitempty"`
	Range float64           `json:",omitempty"`
	Tags  map[string]string `json:",omitempty"`
}
