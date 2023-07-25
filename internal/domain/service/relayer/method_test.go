package relayer_test

import (
	"os"
	"testing"

	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
)



var instance *relayer.RelayerManager

func SetUp() {

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


func Test_Relayer(t *testing.T) {
	t.Skip("")	
}