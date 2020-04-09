package runtime

import (
	"math/rand"
	"time"

	app "github.com/alexwbaule/go-app"
	"github.com/oklog/ulid"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(app app.Application) (*Runtime, error) {

	rt := &Runtime{
		app: app,
	}

	return rt, nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	app app.Application
}
