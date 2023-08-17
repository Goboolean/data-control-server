package kis_test

import (
	"os"
	"testing"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/kis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
	"github.com/Goboolean/shared/pkg/resolver"
	_ "github.com/Goboolean/fetch-server/internal/util/env"
)

var instance ws.Fetcher

var (
	count    int         = 0
	receiver ws.Receiver = mock.NewMockReceiver(func() {
		count++
	})
)

func SetupKis() {
	var err error

	instance, err = kis.New(&resolver.ConfigMap{
		"APPKEY": os.Getenv("KIS_APPKEY"),
		"SECRET": os.Getenv("KIS_SECRET"),
	}, receiver)
	if err != nil {
		panic(err)
	}
}

func TeardownKis() {
	instance.Close()
}

func TestMain(m *testing.M) {
	os.Exit(0)
	SetupKis()
	code := m.Run()
	TeardownKis()
	os.Exit(code)
}

func Test_Constructor(t *testing.T) {

	t.Run("Ping", func(t *testing.T) {
		if err := instance.Ping(); err != nil {
			t.Errorf("Ping() = %v", err)
			return
		}
	})
}
