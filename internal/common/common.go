package common

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// GetULID returns a new ULID
func GetULID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
