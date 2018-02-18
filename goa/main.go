//go:generate goagen bootstrap -d github.com/craftone/gojiko/goa/design

package main

import (
	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/goa/app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	config.Init()
	domain.Init()

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
