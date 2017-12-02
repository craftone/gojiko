package apns

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel)
	_log := logger.WithField("package", "domain")
	Init(_log)

	code := m.Run()
	os.Exit(code)
}
