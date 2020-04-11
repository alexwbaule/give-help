package runtime

import (
	firebase "firebase.google.com/go"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/fireadmin"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	app "github.com/alexwbaule/go-app"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(app app.Application) (*Runtime, error) {

	c := connection.New(&common.DbConfig{
		DBName: app.Config().GetString("database.DBName"),
		Host:   app.Config().GetString("database.Host"),
		Pass:   app.Config().GetString("database.Pass"),
		User:   app.Config().GetString("database.User"),
	})

	rt := &Runtime{
		app:      app,
		fbase:    fireadmin.InitializeAppWithServiceAccount(),
		database: c,
	}

	return rt, nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	app      app.Application
	fbase    *firebase.App
	database *connection.Connection
}

func (rt *Runtime) GetFirebase() *firebase.App {
	return rt.fbase
}

func (rt *Runtime) GetDatabase() *connection.Connection {
	return rt.database
}

func (rt *Runtime) CloseDatabase() {
	rt.database.Close()
}
