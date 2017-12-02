package apns

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

var theRepo = newRepo()

// TheRepo returns the global APN repository in this program.
func TheRepo() *Repo {
	return theRepo
}

func Init(_log *logrus.Entry) error {
	log = _log.WithField("package", "domain/apns")
	log.Info("Initialize APN package")

	return nil
}
