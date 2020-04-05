package common

import "time"

type TransactionStatus string

const (
	Open       TransactionStatus = "open"
	InProgress                   = "in-progess"
	Done                         = "done"
	Cancel                       = "cancel"
)

type TransactionReview struct {
	TransactionId string            `json:",omitempty"`
	CreatedAt     time.time         `json:",omitempty"`
	GiverId       string            `json:",omitempty"`
	TakerId       string            `json:",omitempty"`
	Description   string            `json:",omitempty"`
	Tags          map[string]string `json:",omitempty"`
	GiverReview   *UserReview       `json:",omitempty"`
	TakerReview   *UserReview       `json:",omitempty"`
	Status        TransactionStatus `json:",omitempty"`
}

type UserReview struct {
	Rating  int    `json:",omitempty"`
	Comment string `json:",omitempty"`
}
