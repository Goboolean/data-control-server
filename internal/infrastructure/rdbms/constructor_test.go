package rdbms_test

import (
	"os"
	"testing"

	"github.com/Goboolean/shared/pkg/rdbms"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/joho/godotenv"
)

var (
	db      *rdbms.PSQL
	queries *rdbms.Queries
)

func TestMain(m *testing.M) {

	if err := os.Chdir("../../"); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	SetupPSQL()
	code := m.Run()
	TeardownPSQL()

	os.Exit(code)
}

func SetupPSQL() {

	var err error
	db, err = rdbms.NewDB(&resolver.ConfigMap{
		"HOST":     os.Getenv("PSQL_HOST"),
		"PORT":     os.Getenv("PSQL_PORT"),
		"USER":     os.Getenv("PSQL_USER"),
		"PASSWORD": os.Getenv("PSQL_PASS"),
		"DATABASE": os.Getenv("PSQL_DATABASE"),
	})
	if err != nil {
		panic(err)
	}

	queries = rdbms.NewQueries(db)
}

func TeardownPSQL() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func Test_Constructor(t *testing.T) {
	if err := db.Ping(); err != nil {
		t.Errorf("Ping() failed: %v", err)
	}
}



