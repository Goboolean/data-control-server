package main

import (
	"github.com/Goboolean/fetch-server/cmd/compose/released"
	_ "github.com/Goboolean/fetch-server/internal/util/logger"
	log "github.com/sirupsen/logrus"
)



func main() {
	// An option for released, not now
	//log.Info(released.Run())

	// An option for develop
	log.Info(released.Run())
}