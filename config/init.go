package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	gtpv2c_timeout int
	gtpv2c_retry   int
	mtu            uint16
)

func Init() {
	viper.SetConfigName("gojiko")
	viper.AddConfigPath("$HOME/.gojiko")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	// Set default values
	viper.SetDefault("gtpv2c.timeout", 1000) // msec
	viper.SetDefault("gtpv2c.retry", 2)      // times
	viper.SetDefault("mtu", 1472)            // bytes

	// check and load configs
	err = SetGtpv2cTimeout(viper.GetInt("gtpv2c.timeout"))
	if err != nil {
		log.Fatal(err)
	}
	err = SetGtpv2cRetry(viper.GetInt("gtpv2c.retry"))
	if err != nil {
		log.Fatal(err)
	}
	err = setMTU(uint16(viper.GetInt("mtu")))
	if err != nil {
		log.Fatal(err)
	}

	initApn()
	initSgw()
	initLog()
	initStats()
}

// about Gtpv2cTimeout
func Gtpv2cTimeout() int {
	return gtpv2c_timeout
}

func Gtpv2cTimeoutDuration() time.Duration {
	return time.Duration(gtpv2c_timeout) * time.Millisecond
}

func SetGtpv2cTimeout(msec int) error {
	gtpv2c_timeout = msec
	return nil
}

// about Gtpv2cRetry
func Gtpv2cRetry() int {
	return gtpv2c_retry
}

func SetGtpv2cRetry(val int) error {
	if val < 0 {
		return errors.New("Invalid gtpv2c.retry")
	}
	gtpv2c_retry = val
	return nil
}

// about MTU
func MTU() uint16 {
	return mtu
}

func setMTU(val uint16) error {
	if val < 600 {
		return fmt.Errorf("Too short MTU : %d", val)
	}
	mtu = val
	return nil
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
