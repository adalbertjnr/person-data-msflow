package types

import (
	"github.com/sirupsen/logrus"
)

func LoggerJSONClient() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}
