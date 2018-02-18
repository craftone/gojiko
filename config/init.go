package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	gtpv2c_timeout int
	gtpv2c_retry   int
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
	viper.SetDefault("gtpv2c.timeout", 1000) // msec
	viper.SetDefault("gtpv2c.retry", 2)      // times

	// check and load configs
	gtpv2c_timeout = viper.GetInt("gtpv2c.timeout")
	gtpv2c_retry = viper.GetInt("gtpv2c.retry")

	initApn()
	initSgw()
}

// about Gtpv2cTimeout
func Gtpv2cTimeout() int {
	return gtpv2c_timeout
}

func Gtpv2cTimeoutDuration() time.Duration {
	return time.Duration(gtpv2c_timeout) * time.Millisecond
}

func SetGtpv2cTimeout(msec int) {
	gtpv2c_timeout = msec
}

// about Gtpv2cRetry
func Gtpv2cRetry() int {
	return gtpv2c_retry
}

func SetGtpv2cRetry(val int) {
	gtpv2c_retry = val
}

func getStringFromIfMap(ifMap map[string]interface{}, genre, key string) (string, error) {
	strInt, exists := ifMap[key]
	if !exists {
		return "", fmt.Errorf("Invalid %s config : No %s", genre, key)
	}
	str, ok := strInt.(string)
	if !ok {
		return "", fmt.Errorf("Invalid %s config : The %s is not a string : %#v", genre, key, strInt)
	}
	return str, nil
}
