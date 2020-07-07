package runtime

import (
	"log"

	firebase "firebase.google.com/go"
	metrics "git.corp.c6bank.com/c6libs/go-c6-metrics"
	cacheConn "github.com/alexwbaule/give-help/v2/internal/cache/connection"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/fireadmin"
	dbConn "github.com/alexwbaule/give-help/v2/internal/storage/connection"
	app "github.com/alexwbaule/go-app"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(app app.Application) (*Runtime, error) {
	db := dbConn.New(&common.DbConfig{
		DBName: app.Config().GetString("database.DBName"),
		Host:   app.Config().GetString("database.Host"),
		Pass:   app.Config().GetString("database.Pass"),
		User:   app.Config().GetString("database.User"),
	})

	es, err := cacheConn.New(&common.CacheConfig{
		Addresses: app.Config().GetStringSlice("es.Addresses"),
	})

	if err != nil {
		log.Printf("fail to connect on cache: %s\n", err)
		return nil, err
	}

	firebaseAccountKeyPath := app.Config().GetString("firebase.AccountKey")

	if len(firebaseAccountKeyPath) == 0 {
		firebaseAccountKeyPath = `etc/serviceAccountKey.json`
	}

	rt := &Runtime{
		app:      app,
		fbase:    fireadmin.InitializeAppWithServiceAccount(firebaseAccountKeyPath),
		database: db,
		cache:    es,
		Metrics:  metrics.NewResource(metrics.Config{}),
	}

	return rt, err
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	app      app.Application
	fbase    *firebase.App
	database *dbConn.Connection
	Metrics  *metrics.Resources
	cache    *cacheConn.Connection
}

func (rt *Runtime) GetCache() *cacheConn.Connection {
	return rt.cache
}

func (rt *Runtime) GetFirebase() *firebase.App {
	return rt.fbase
}

func (rt *Runtime) GetDatabase() *dbConn.Connection {
	return rt.database
}

func (rt *Runtime) CloseDatabase() {
	rt.database.Close()
}

func (rt *Runtime) GetMetricProcessor() *metrics.Processor {
	return rt.metricProcessor
}
