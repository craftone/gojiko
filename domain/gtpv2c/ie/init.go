package ie

import (
	"github.com/craftone/gojiko/applog"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func Init() {
	log = applog.NewLogger("gtpv2c/ie")
	log.Infof("Package initialized")
}
