package domain

import (
	"net"

	"github.com/craftone/gojiko/applog"
	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain/apns"
	"github.com/craftone/gojiko/domain/gtpv2c"
	"github.com/craftone/gojiko/domain/stats"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

var theSgwCtrlRepo *sgwCtrlRepo
var defaultSgwCtrlAddr = net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: GtpControlPort}

func Init() error {
	log = applog.NewLogEntry("domain")
	log.Info("Initialize domain package")

	apns.Init()
	gtpv2c.Init()
	stats.Init()

	theSgwCtrlRepo = newSgwCtrlRepo()
	for _, sgw := range config.GetSGWs() {
		sgwCtrlAddr := net.UDPAddr{IP: sgw.IP, Port: GtpControlPort}
		sgwCtrl, err := newSgwCtrl(sgwCtrlAddr, GtpUserPort, 0)
		if err != nil {
			return err
		}
		theSgwCtrlRepo.AddCtrl(sgwCtrl)
	}

	return nil
}

func TheSgwCtrlRepo() *sgwCtrlRepo {
	return theSgwCtrlRepo
}
