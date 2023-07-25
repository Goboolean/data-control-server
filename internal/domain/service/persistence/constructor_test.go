package persistence_test

import (
	"os"
	"testing"

	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
)





var instance *persistence.PersistenceManager


func SetUp() {
	//instance = persistence.New()
}


func TearDown() {
	instance.Close()
}


func TestMain(m *testing.M) {

	SetUp()
	code := m.Run()
	TearDown()
	os.Exit(code)

}


