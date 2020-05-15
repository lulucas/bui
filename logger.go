package bui

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Logger() *logrus.Logger {
	return log
}

func init() {
	log.SetLevel(logrus.DebugLevel)
}
