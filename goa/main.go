//go:generate goagen bootstrap -d github.com/craftone/gojiko/goa/design

package main

import (
	"log"

	"github.com/craftone/gojiko/applog"
	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/goa/app"
	"github.com/goadesign/goa"
	goalogrus "github.com/goadesign/goa/logging/logrus"
	"github.com/goadesign/goa/middleware"
)

func main() {
	config.Init()
	err := domain.Init()
	if err != nil {
		log.Panic(err)
	}

	// Create service
	service := goa.New("gojiko api")

	// set log adapter
	logger := applog.NewLogger("goa")
	service.WithLogger(goalogrus.New(logger))

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "gtpsession" controller
	c1 := NewGtpsessionController(service)
	app.MountGtpsessionController(service, c1)

	// Mount "udpEchoFlowByIMSIandEBI" controller
	c2 := NewUDPEchoFlowByIMSIandEBIController(service)
	app.MountUDPEchoFlowByIMSIandEBIController(service, c2)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}
