package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)


// !!! Import this package at main.
// _ "github.com/Goboolean/fetch-server/internal/util/logger"
func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})
}