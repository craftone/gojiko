package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var logLevel string

func initLog() {
	viper.SetDefault("log.level", "info")
	logLevelStr := viper.GetString("log.level")
	logLevel = strings.ToLower(logLevelStr)
	switch logLevel {
	case "debug", "info", "warn", "error":
	default:
		log.Fatalf("Invalid log level is configured: %s", logLevelStr)
	}
}

func LogLevel(pkgName string) string {
	// pkgName is for package stratum log level, but not yet implemented
	return logLevel
}
