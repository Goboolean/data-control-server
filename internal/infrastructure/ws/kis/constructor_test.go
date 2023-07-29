package kis_test

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/kis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
	"github.com/Goboolean/fetch-server/internal/util/env"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/joho/godotenv"
)

var instance ws.Fetcher

var (
	count    int         = 0
	receiver ws.Receiver = mock.NewMockReceiver(func() {
		count++
	})
)

func SetupKis() {

	if err := os.Chdir(env.Root); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	instance = kis.New(&resolver.ConfigMap{
		"APPKEY": os.Getenv("KIS_APPKEY"),
		"SECRET": os.Getenv("KIS_SECRET"),
	}, context.Background(), receiver)
}

func TeardownKis() {
	instance.Close()
}

func TestMain(m *testing.M) {

	if err := os.Chdir(env.Root); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

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
