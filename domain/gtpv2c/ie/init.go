package ie

import (
	"github.com/craftone/gojiko/applog"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func Init() {
	log = applog.NewLogEntry("gtpv2c/ie")
}
