package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func InitLogger() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

func GetLogger() *logrus.Logger {
	return logger
}
