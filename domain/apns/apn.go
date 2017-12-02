package apns

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// Apn is a Value-Object
type Apn struct {
	networkID  string
	operatorID string
	fullString string
	ips        []net.IP
	curIdx     int
	mtx        sync.Mutex
}

func NewApn(networkID, mcc, mnc string, ips []net.IP) (*Apn, error) {
	if err := checkApnNetworkID(networkID); err != nil {
		return nil, err
	}
	if err := checkApnOperatorID(mcc, mnc); err != nil {
		return nil, err
	}
	if len(mnc) == 2 {
		mnc = "0" + mnc
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("One or more IPs should be specified")
	}

	operatorID := fmt.Sprintf("mnc%s.mcc%s.gprs", mnc, mcc)

	return &Apn{
		networkID:  networkID,
		operatorID: operatorID,
		fullString: strings.ToLower(networkID) + "." + operatorID,
		ips:        ips,
		curIdx:     -1, // before first use, curIdx will be incremented.
	}, nil
}

func (a *Apn) FullString() string {
	return a.fullString
}

func (a *Apn) GetIP() net.IP {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.curIdx++
	if a.curIdx >= len(a.ips) {
		a.curIdx = 0
	}
	return a.ips[a.curIdx]
}

func checkApnNetworkID(networkID string) error {
	if len(networkID) > 63 {
		return fmt.Errorf("Too long string for APN : %s", networkID)
	}

	u, err := url.Parse("apn://" + networkID)
	if err != nil {
		return err
	}

	if len(u.Hostname()) == 0 {
		return fmt.Errorf("Invalid APN Network Identifier : %s", networkID)
	}

	lnetworkID := strings.ToLower(networkID)

	if strings.ContainsRune(lnetworkID, '*') {
		return fmt.Errorf("APN Network Identifier should not contain * : %s", networkID)
	}

	if match, _ := regexp.MatchString("^([rl]ac|sgsn|rnc)", lnetworkID); match {
		return fmt.Errorf("APN Network Identifier should not start with that: %s", networkID)
	}

	if strings.HasSuffix(lnetworkID, ".gprs") {
		return fmt.Errorf("APN Network Identifier should not ends with .gprs : %s", networkID)
	}

	return nil
}

func checkApnOperatorID(mcc, mnc string) error {
	if match, _ := regexp.MatchString("^[0-9]{3}$", mcc); !match {
		return fmt.Errorf("MCC should be a 3-digits string : %s", mcc)
	}
	if match, _ := regexp.MatchString("^[0-9]{2,3}$", mnc); !match {
		return fmt.Errorf("MNC should be a 3-digits or 2-digits string : %s", mnc)
	}
	return nil
}
