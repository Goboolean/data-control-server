package mock_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
	"github.com/joho/godotenv"
)




var (
	instance ws.Fetcher
	receiver ws.Receiver
)


func SetupMock() {
	instance = mock.New(context.Background(), 1 * time.Millisecond, receiver)
}

func TeardownMock() {
	instance.Close()
}


func TestMain(m *testing.M) {

	if err := os.Chdir("../../../"); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	SetupMock()
	code := m.Run()
	TeardownMock()

	os.Exit(code)
}


func Test_Constructor(t *testing.T) {

	if err := instance.Ping(); err != nil {
		t.Errorf("Ping() = %v", err)
	}
}

