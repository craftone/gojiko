package config

import (
	"log"
	"net"

	"github.com/spf13/viper"
)

type Sgw struct {
	Name     string
	IP       net.IP
	Recovery byte
}

var (
	sgws []Sgw
)

func loadSgwConf() {
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
			log.Panicf("Invalid SGW config : type assertion error %#v", sgwMapIf)
		}
		name, err := getStringFromIfMap(sgwMap, "SGW", "name", "")
		if err != nil {
			log.Panicf(err.Error())
		}

		ipStr, err := getStringFromIfMap(sgwMap, "SGW", "ip", "")
		if err != nil || ipStr == "" {
			log.Panicf(err.Error())
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			log.Panicf("Invalid SGW config : %s is not a valid IP", ipStr)
		}

		rec, err := getInt64FromIfMap(sgwMap, "SGW", "recovery", 0)
		if err != nil {
			log.Panicf(err.Error())
		}
		if rec < 0 || rec > 0xff {
			log.Panicf("recovery should be a number from 0 to 255")
		}

		sgws = append(sgws, Sgw{name, ip, byte(rec)})
	}
}

func GetSGWs() []Sgw {
	return sgws
}
