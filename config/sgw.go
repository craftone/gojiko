package config

import (
	"log"
	"net"

	"github.com/spf13/viper"
)

type Sgw struct {
	Name string
	IP   net.IP
}

var (
	sgws []Sgw
)

func initSgw() {
	if !viper.IsSet("sgw") {
		log.Panic("Invalid SGW config : No SGW config")
	}
	sgwMapIfs, ok := viper.Get("sgw").([]interface{})
	if !ok {
		log.Panicf("Invalid SGW config : type assertion error %#v", viper.Get("sgw"))
	}
	for _, sgwMapIf := range sgwMapIfs {
		sgwMap, ok := sgwMapIf.(map[string]interface{})
		if !ok {
			log.Panicf("Invalid SGW config : type asseertion error %#v", sgwMapIf)
		}
		name, err := getStringFromIfMap(sgwMap, "SGW", "name")
		if err != nil {
			log.Panicf(err.Error())
		}
		ipStr, err := getStringFromIfMap(sgwMap, "SGW", "ip")
		if err != nil {
			log.Panicf(err.Error())
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			log.Panicf("Invalid SGW config : %s is not a valid IP", ipStr)
		}
		sgws = append(sgws, Sgw{name, ip})
	}
}

func GetSGWs() []Sgw {
	return sgws
}
