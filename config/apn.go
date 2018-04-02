package config

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/viper"
)

type Apn struct {
	Host string
	Mcc  string
	Mnc  string
	IPs  []net.IP
}

var (
	apns []Apn
)

func loadApnConf() {
	if !viper.IsSet("apn") {
		log.Panic("Invalid APN config : No APN config")
	}
	apnMapIfs, ok := viper.Get("apn").([]interface{})
	if !ok {
		log.Panicf("Invalid APN config : type assertion error %#v", viper.Get("apn"))
	}
	for _, apnMapIf := range apnMapIfs {
		apnMap, ok := apnMapIf.(map[string]interface{})
		if !ok {
			log.Panicf("Invalid APN config : type asseertion error %#v", apnMapIf)
		}
		apn, err := newApn(apnMap)
		if err != nil {
			log.Panic(err)
		}
		apns = append(apns, apn)
	}
}

func newApn(apnMap map[string]interface{}) (Apn, error) {
	name, err := getStringFromIfMap(apnMap, "APN", "host", "")
	if err != nil {
		return Apn{}, err
	}
	if name == "" {
		return Apn{}, fmt.Errorf("No name in APN config")
	}
	mcc, err := getStringFromIfMap(apnMap, "APN", "mcc", "")
	if err != nil {
		return Apn{}, err
	}
	if mcc == "" {
		return Apn{}, fmt.Errorf("No mcc in APN config")
	}
	mnc, err := getStringFromIfMap(apnMap, "APN", "mnc", "")
	if err != nil || mnc == "" {
		return Apn{}, err
	}
	if mnc == "" {
		return Apn{}, fmt.Errorf("No mnc in APN config")
	}
	ips, err := getIPsFromApnMap(apnMap, "ips")
	if err != nil {
		return Apn{}, err
	}
	return Apn{name, mcc, mnc, ips}, nil
}

func getIPsFromApnMap(apnMap map[string]interface{}, key string) ([]net.IP, error) {
	ips := []net.IP{}
	ipsIf, exists := apnMap[key]
	if !exists {
		return ips, fmt.Errorf("Invalid APN config : No %s", key)
	}
	ipStrsIf, ok := ipsIf.([]interface{})
	if !ok {
		return ips, fmt.Errorf("Invalid APN config : The %s is not an array of string : %#v", key, ipStrsIf)
	}

	for _, ipStrIf := range ipStrsIf {
		ipStr, ok := ipStrIf.(string)
		if !ok {
			return ips, fmt.Errorf("Invalid APN config : The %s is not a string : %#v", key, ipStrIf)
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return ips, fmt.Errorf("Invalid APN config : %s is not a valid IP", ipStr)
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

func GetAPNs() []Apn {
	return apns
}
