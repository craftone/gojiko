//go:generate goagen bootstrap -d github.com/craftone/gojiko/goa/design

package main

import (
	"net"

	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/domain/apns"
	"github.com/craftone/gojiko/goa/app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

var pgwIP = net.IPv4(127, 1, 1, 1)
var apn, _ = apns.NewApn("example.com", "440", "10", []net.IP{pgwIP})
var defaultSgwCtrlAddr = net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: domain.GtpControlPort}

func main() {
	config.Init()
	domain.Init()
	// default PGW
	apns.TheRepo().Post(apn)

	// Create service
	service := goa.New("gojiko api")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "gtpsession" controller
	c := NewGtpsessionController(service)
	app.MountGtpsessionController(service, c)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}
