package runtime

import (
	"encoding/json"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/alexwbaule/give-help/v2/internal/fireadmin"
	"github.com/alexwbaule/give-help/v2/internal/storage/connection"
	app "github.com/alexwbaule/go-app"
	"github.com/elastic/go-elasticsearch"
	"github.com/rafaelfino/metrics"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(app app.Application) (*Runtime, error) {
	c := connection.New(&common.DbConfig{
		DBName: app.Config().GetString("database.DBName"),
		Host:   app.Config().GetString("database.Host"),
		Pass:   app.Config().GetString("database.Pass"),
		User:   app.Config().GetString("database.User"),
	})

	metricsInterval := app.Config().GetString("metrics.Interval")

	interval, err := time.ParseDuration(metricsInterval)

	if err != nil {
		interval = time.Minute * 10
	}

	rt := &Runtime{
		app:             app,
		fbase:           fireadmin.InitializeAppWithServiceAccount(app.Config().GetString("firebase.AccountKey")),
		database:        c,
		metricProcessor: metrics.NewMetricProcessor(interval, LogExport),
	}

	return rt, nil
}

func LogExport(data *metrics.MetricData) error {
	raw, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		log.Printf("fail to marshall metrics: %s\n", err)
	}

	log.Printf("Metrics: %s\n", string(raw))

	return nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	app             app.Application
	fbase           *firebase.App
	database        *connection.Connection
	metricProcessor *metrics.Processor
}

func (rt *Runtime) GetElasticSearchConfig() *elasticsearch.Config {
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

func (rt *Runtime) GetMetricProcessor() *metrics.Processor {
	return rt.metricProcessor
}
