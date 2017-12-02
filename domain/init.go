package domain

import (
	"net"
	"os"

	"github.com/craftone/gojiko/domain/apns"
	"github.com/sirupsen/logrus"
)

var log = newLogger()

var theGtpSessionRepo *gtpSessionRepo
var theSgwCtrlRepo *sgwCtrlRepo
var defaultSgwCtrlAddr = net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: GtpControlPort}

func Init() error {
	log.Info("Initialize domain package")

	apns.Init(log)
	sgwCtrl, err := newSgwCtrl(defaultSgwCtrlAddr, GtpUserPort, 0)
	if err != nil {
		return err
	}
	theSgwCtrlRepo = newSgwCtrlRepo()
	theSgwCtrlRepo.AddCtrl(sgwCtrl)

	theGtpSessionRepo = newGtpSessionRepo()

	return nil
}

func newLogger() *logrus.Entry {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel)
	return logger.WithField("package", "domain")
}

func TheSgwCtrlRepo() *sgwCtrlRepo {
	return theSgwCtrlRepo
}
