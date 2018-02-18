package apns

import (
	"github.com/craftone/gojiko/config"
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

	for _, capn := range config.GetAPNs() {
		apn, err := NewApn(capn.Host, capn.Mcc, capn.Mnc, capn.IPs)
		if err != nil {
			log.Fatalf("Invalid APN config : %#v", capn)
		}
		err = theRepo.Post(apn)
		if err != nil {
			log.Fatal("APN Post error : %v", err)
		}
	}

	return nil
}
