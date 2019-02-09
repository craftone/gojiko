package domain

import (
	"net"
	"os"
	"testing"

	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain/apns"
)

var pgwIP = net.IPv4(127, 1, 1, 1).To4()
var pgwDataIP = net.IPv4(127, 1, 1, 2).To4()
var apn, _ = apns.NewApn("example.com", "440", "10", []net.IP{pgwIP})

func TestMain(m *testing.M) {
	config.Init()
	Init()

	// default PGW
	apns.TheRepo().Post(apn)

	code := m.Run()
	os.Exit(code)
}
