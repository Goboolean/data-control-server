package buycycle_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/buycycle"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
	_ "github.com/Goboolean/fetch-server/internal/util/env"
	"github.com/Goboolean/fetch-server/internal/util/withintime"
	"github.com/Goboolean/shared/pkg/resolver"
)



var instance ws.Fetcher

var (
	count int = 0
	receiver ws.Receiver = mock.NewMockReceiver(func() {
		count++
	})
)





func SetupBuycycle() {
	var err error

	instance, err = buycycle.New(&resolver.ConfigMap{
		"HOST": os.Getenv("BUYCYCLE_HOST"),
		"PORT": os.Getenv("BUYCYCLE_PORT"),
		"PATH": os.Getenv("BUYCYCLE_PATH"),
	}, receiver)
	if err != nil {
		panic(err)
	}
}

func TeardownBuycycle() {
	instance.Close()
}

func TestMain(m *testing.M) {
	os.Exit(0)
	SetupBuycycle()
	code := m.Run()
	TeardownBuycycle()
	os.Exit(code)
	os.Exit(m.Run())
}




func Test_Constructor(t *testing.T) {
	t.Skip()
}


var (
	once sync.Once
	withinDurationChecker *withintime.WithinDurationChecker
	isMarketOnCache bool
)

// Struct withinDurationChecker is initialized with information of the Korea stock market.
// Value isMarketOnCache is cached at the time of first call, therefore inconsistency beween tests may not occur.
// U can get the value by calling IsMarketOn() function.
func isMarketOn() bool {

	once.Do(func() {
		var err error

		withinDurationChecker, err = withintime.New(&withintime.Option{
			Location: "Asia/Seoul",
			Inclusion: withintime.ConditionList{
				time.Monday:    &withintime.Condition{StartTime: "09:00", EndTime: "15:30"},
				time.Tuesday:   &withintime.Condition{StartTime: "09:00", EndTime: "15:30"},
				time.Wednesday: &withintime.Condition{StartTime: "09:00", EndTime: "15:30"},
				time.Thursday:  &withintime.Condition{StartTime: "09:00", EndTime: "15:30"},
				time.Friday:    &withintime.Condition{StartTime: "09:00", EndTime: "15:30"},
			},
			Exclusion: withintime.ConditionList{
				// TODO: add holiday
				},
			}, nil)

		if err != nil {
			panic(err)
		}

		isMarketOnCache = withinDurationChecker.IsTimeNowWithinDuration()
	})
	
	return isMarketOnCache
}