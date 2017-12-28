package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("gojiko")
	viper.AddConfigPath("$HOME/.gojiko")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	// Set default values
	viper.SetDefault("gtpv2c.timeout", 3000) // msec
	viper.SetDefault("gtpv2c.retry", 3)      // times
}

// about Gtpv2cTimeout
func Gtpv2cTimeout() int {
	return viper.GetInt("gtpv2c.timeout")
}

func Gtpv2cTimeoutDuration() time.Duration {
	return time.Duration(viper.GetInt("gtpv2c.timeout")) * time.Millisecond
}

func SetGtpv2cTimeout(msec int) {
	viper.Set("gtpv2c.timeout", msec)
}

// about Gtpv2cRetry
func Gtpv2cRetry() int {
	return viper.GetInt("gtpv2c.retry")
}

func SetGtpv2cRetry(val int) {
	viper.Set("gtpv2c.retry", val)
}
