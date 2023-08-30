package main

import (
	"github.com/Goboolean/fetch-server.v1/cmd/compose/released"
	_ "github.com/Goboolean/fetch-server.v1/internal/util/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	// An option for released, not now
	//log.Info(released.Run())

	// An option for develop
	log.Info(released.Run())
}
