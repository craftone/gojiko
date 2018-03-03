package applog

import (
	"log"
	"os"

	"github.com/craftone/gojiko/config"
	"github.com/sirupsen/logrus"
)

func Init() {
	// Nothing yet
	// TODO: open thread-safe log file
	//       ref https://github.com/sirupsen/logrus/issues/391
}

func NewLogger(pkgName string) *logrus.Entry {
	logger := logrus.New()

	logger.Out = os.Stdout

	logLevelStr := config.LogLevel(pkgName)
	logger.WithField("package", pkgName).Infof("A logger created at %s level", logLevelStr)

	switch logLevelStr {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		log.Panicf("Invalid logLevel : %s", logLevelStr)
	}
	myLog := logger.WithField("package", pkgName)
	return myLog
}
