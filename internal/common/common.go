package common

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/oklog/ulid"
)

// GetULID returns a new ULID
func GetULID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

const DegreeInKM float64 = 111

//CalculeRange with a target area, calcule region with range in km
func CalculeRange(area *models.Location) (float64, float64, float64, float64, error) {
	var latN float64
	var latS float64
	var lonW float64
	var lonE float64

	var err error

	if area == nil || area.Lat == nil || area.Lon == nil {
		err = fmt.Errorf("cannot calculate nil model.Area")
		return latN, latS, lonW, lonE, err
	}

	latN = round(*area.Lat + area.Distance/DegreeInKM)
	latS = round(*area.Lat - area.Distance/DegreeInKM)
	lonW = round(*area.Lon - math.Cos(*area.Lat)*area.Distance/DegreeInKM)
	lonE = round(*area.Lon + math.Cos(*area.Lat)*area.Distance/DegreeInKM)

	return latN, latS, lonW, lonE, err
}

func round(x float64) float64 {
	return math.Round(x*10000) / 10000
}

type CacheConfig struct {
	Addresses []string
}

//Config base connection config struct
type DbConfig struct {
	Host   string
	User   string
	Pass   string
	DBName string
}

func NormalizeTagArray(arr []string) []string {
	m := map[string]interface{}{}

	for _, a := range arr {
		m[strings.ToLower(a)] = nil
	}

	ret := []string{}

	for k := range m {
		ret = append(ret, k)
	}

	return ret
}
