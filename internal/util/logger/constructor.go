package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)



func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})
}