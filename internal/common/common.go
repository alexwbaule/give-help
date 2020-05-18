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
	var longW float64
	var longE float64

	var err error

	if area == nil || area.Lat == nil || area.Long == nil {
		err = fmt.Errorf("cannot calculate nil model.Area")
		return latN, latS, longW, longE, err
	}

	latN = round(*area.Lat + area.Range/DegreeInKM)
	latS = round(*area.Lat - area.Range/DegreeInKM)
	longW = round(*area.Long - math.Cos(*area.Lat)*area.Range/DegreeInKM)
	longE = round(*area.Long + math.Cos(*area.Lat)*area.Range/DegreeInKM)

	return latN, latS, longW, longE, err
}

func round(x float64) float64 {
	return math.Round(x*10000) / 10000
}

type Config struct {
	Db *DbConfig
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
