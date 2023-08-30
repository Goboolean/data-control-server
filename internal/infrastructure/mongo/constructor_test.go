package mongo_test

import (
	"os"
	"testing"

	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/joho/godotenv"
)

// Primary test for mongoDB should contain the following:
// 1. Test if server is running (ping test)
// 2. Test insertion without transaction
// 3. Test insertion with transaction

var (
	instance *mongo.DB
	queries  *mongo.Queries
)

func SetupMongo() {

	var err error
	instance, err = mongo.NewDB(&resolver.ConfigMap{
		"HOST":     os.Getenv("MONGO_HOST"),
		"USER":     os.Getenv("MONGO_USER"),
		"PORT":     os.Getenv("MONGO_PORT"),
		"PASSWORD": os.Getenv("MONGO_PASS"),
		"DATABASE": os.Getenv("MONGO_DATABASE"),
	})
	if err != nil {
		panic(err)
	}
	queries = mongo.New(instance)
}

func TeardownMongo() {
	instance.Close()
}

func TestMain(m *testing.M) {

	if err := os.Chdir("../../"); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	SetupMongo()
	code := m.Run()
	TeardownMongo()

	os.Exit(code)
}

func TestConstructor(t *testing.T) {
	if err := instance.Ping(); err != nil {
		t.Errorf("Ping() failed: %v", err)
	}
}
