package gtpv2c

import (
	"github.com/craftone/gojiko/applog"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func Init() {
	ie.Init()
	log = applog.NewLogger("gtpv2c")
}
