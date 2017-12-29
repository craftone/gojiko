package domain

import (
	"net"

	"github.com/sirupsen/logrus"
)

type SgwData struct {
	*absSPgw
}

func newSgwData(addr net.UDPAddr, recovery byte, sgwCtrl *SgwCtrl) (*SgwData, error) {
	myLog := log.WithFields(logrus.Fields{
		"addr":     addr,
		"recovery": recovery,
	})
	absSPgw, err := newAbsSPgw(addr, recovery, sgwCtrl)
	if err != nil {
		return nil, err
	}
	myLog.Debug("A new SGW Data is created")
	return &SgwData{absSPgw}, nil
}