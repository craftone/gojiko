package main

import (
	"github.com/craftone/gojiko/goa/app"
	"github.com/goadesign/goa"
)

// GtpsessionController implements the gtpsession resource.
type GtpsessionController struct {
	*goa.Controller
}

// NewGtpsessionController creates a gtpsession controller.
func NewGtpsessionController(service *goa.Service) *GtpsessionController {
	return &GtpsessionController{Controller: service.NewController("GtpsessionController")}
}

// Create runs the create action.
func (c *GtpsessionController) Create(ctx *app.CreateGtpsessionContext) error {
	// GtpsessionController_Create: start_implement

	// Put your logic here

	res := &app.Gtpsession{}
	return ctx.OK(res)
	// GtpsessionController_Create: end_implement
}

// ShowByID runs the showByID action.
func (c *GtpsessionController) ShowByID(ctx *app.ShowByIDGtpsessionContext) error {
	// GtpsessionController_ShowByID: start_implement

	// Put your logic here

	res := &app.Gtpsession{}
	return ctx.OK(res)
	// GtpsessionController_ShowByID: end_implement
}

// ShowByIMSIandEBI runs the showByIMSIandEBI action.
func (c *GtpsessionController) ShowByIMSIandEBI(ctx *app.ShowByIMSIandEBIGtpsessionContext) error {
	// GtpsessionController_ShowByIMSIandEBI: start_implement

	// Put your logic here

	res := &app.Gtpsession{}
	return ctx.OK(res)
	// GtpsessionController_ShowByIMSIandEBI: end_implement
}
