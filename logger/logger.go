package logger

import (
	"github.com/Sirupsen/logrus"
	"openhab2awsiot/config"
)

func Init(config *config.Config) *logrus.Logger {
	log := logrus.New()
	if config.Debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
	if config.Console {
		log.Formatter = &logrus.TextFormatter{}
	}
	return log
}
