package polygon_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/ws/polygon"
	_ "github.com/Goboolean/fetch-server.v1/internal/util/env"
	"github.com/Goboolean/fetch-server.v1/internal/util/withintime"
	"github.com/Goboolean/shared/pkg/resolver"
)

var instance ws.Fetcher

func SetupPolygon() {
	var err error

	instance, err = polygon.New(&resolver.ConfigMap{
		"KEY": os.Getenv("POLYGON_API_KEY"),
	}, receiver)
	if err != nil {
		panic(err)
	}
}

func TeardownPolygon() {
	instance.Close()
}

func TestMain(m *testing.M) {
	os.Exit(0)
	SetupPolygon()
	code := m.Run()
	TeardownPolygon()
	os.Exit(code)
}

func Test_Constructor(t *testing.T) {

	t.Skip("Skip this test, as polygon api key is expired.")

	t.Run("Ping", func(t *testing.T) {
		if err := instance.Ping(); err != nil {
			t.Errorf("Ping() = %v", err)
			return
		}
	})
}

var (
	once                  sync.Once
	withinDurationChecker *withintime.WithinDurationChecker
	isMarketOnCache       bool
)

// Struct withinDurationChecker is initialized with information of the USA stock market.
// Value isMarketOnCache is cached at the time of first call, therefore inconsistency beween tests may not occur.
// U can get the value by calling IsMarketOn() function.
func isMarketOn() bool {

	once.Do(func() {
		var err error

		withinDurationChecker, err = withintime.New(&withintime.Option{
			Location: "America/New_York",
			Inclusion: withintime.ConditionList{
				time.Monday:    &withintime.Condition{StartTime: "09:30", EndTime: "16:00"},
				time.Tuesday:   &withintime.Condition{StartTime: "09:30", EndTime: "16:00"},
				time.Wednesday: &withintime.Condition{StartTime: "09:30", EndTime: "16:00"},
				time.Thursday:  &withintime.Condition{StartTime: "09:30", EndTime: "16:00"},
				time.Friday:    &withintime.Condition{StartTime: "09:30", EndTime: "16:00"},
			},
			Exclusion: withintime.ConditionList{
				// TODO: add holiday
			},
		}, nil)

		if err != nil {
			panic(err)
		}
	})

	return isMarketOnCache
}
