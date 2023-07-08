package polygon_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/polygon"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/util/withintime"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/joho/godotenv"
)



var (
	instance ws.Fetcher
	receiver ws.Receiver
)


func SetupPolygon() {
	instance = polygon.New(&resolver.ConfigMap{
		"KEY":  os.Getenv("POLYGON_API_KEY"),
	}, receiver)
}

func TeardownPolygon() {
	instance.Close()
}


func TestMain(m *testing.M) {

	if err := os.Chdir("../../../"); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	SetupPolygon()
	code := m.Run()
	TeardownPolygon()

	os.Exit(code)
}


func Test_Constructor(t *testing.T) {
	if err := instance.Ping(); err != nil {
		t.Errorf("Ping() = %v", err)
	}
}



var (
	once sync.Once
	withinDurationChecker *withintime.WithinDurationChecker
	isMarketOnCache bool
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
				time.Monday:    &withintime.Condition{StartTime: "09:30", EndTime:   "16:00",},
				time.Tuesday:   &withintime.Condition{StartTime: "09:30", EndTime:   "16:00",},
				time.Wednesday: &withintime.Condition{StartTime: "09:30", EndTime:   "16:00",},
				time.Thursday:  &withintime.Condition{StartTime: "09:30", EndTime:   "16:00",},
				time.Friday:    &withintime.Condition{StartTime: "09:30", EndTime:   "16:00",},
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