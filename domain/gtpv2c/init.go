package gtpv2c

import (
	"github.com/craftone/gojiko/applog"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func Init() {
	ie.Init()
	log = applog.NewLogEntry("gtpv2c")
	log.Info("Initialize GTPv2-C package")
}
