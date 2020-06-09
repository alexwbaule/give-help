package runtime

import (
	"encoding/json"
	"log"
	"time"

	firebase "firebase.google.com/go"
	cacheConn "github.com/alexwbaule/give-help/v2/internal/cache/connection"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/fireadmin"
	dbConn "github.com/alexwbaule/give-help/v2/internal/storage/connection"
	app "github.com/alexwbaule/go-app"
	"github.com/rafaelfino/metrics"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(app app.Application) (*Runtime, error) {
	db := dbConn.New(&common.DbConfig{
		DBName: app.Config().GetString("database.DBName"),
		Host:   app.Config().GetString("database.Host"),
		Pass:   app.Config().GetString("database.Pass"),
		User:   app.Config().GetString("database.User"),
	})

	es, err := cacheConn.New(&cacheConn.Config{
		Addresses: app.Config().GetStringSlice("es.Addresses"),
	})

	if err != nil {
		log.Printf("fail to connect on cache: %s\n", err)
		return nil, err
	}

	metricsInterval := app.Config().GetString("metrics.Interval")

	interval, err := time.ParseDuration(metricsInterval)

	if err != nil {
		interval = time.Minute * 10
	}

	firebaseAccountKeyPath := app.Config().GetString("firebase.AccountKey")

	if len(firebaseAccountKeyPath) == 0 {
		firebaseAccountKeyPath = `etc/serviceAccountKey.json`
	}

	rt := &Runtime{
		app:             app,
		fbase:           fireadmin.InitializeAppWithServiceAccount(firebaseAccountKeyPath),
		database:        db,
		cache:           es,
		metricProcessor: metrics.NewMetricProcessor(interval, LogExport),
	}

	return rt, err
}

func LogExport(data *metrics.MetricData) error {
	raw, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		log.Printf("fail to marshal metrics: %s\n", err)
	} else {
		log.Printf("Metrics: %s\n", string(raw))
	}

	return nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	app             app.Application
	fbase           *firebase.App
	database        *dbConn.Connection
	metricProcessor *metrics.Processor
	cache           *cacheConn.Connection
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
