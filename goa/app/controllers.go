// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "gojiko api": Application Controllers
//
// Command:
// $ goagen
// --design=github.com/craftone/gojiko/goa/design
// --out=$(GOPATH)/src/github.com/craftone/gojiko/goa
// --regen=true
// --version=v1.3.0

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Encoder.Register(goa.NewGobEncoder, "application/gob", "application/x-gob")
	service.Encoder.Register(goa.NewXMLEncoder, "application/xml")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")
	service.Decoder.Register(goa.NewGobDecoder, "application/gob", "application/x-gob")
	service.Decoder.Register(goa.NewXMLDecoder, "application/xml")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// GtpsessionController is the controller interface for the Gtpsession actions.
type GtpsessionController interface {
	goa.Muxer
	Create(*CreateGtpsessionContext) error
}

// MountGtpsessionController "mounts" a Gtpsession resource controller on the given service.
func MountGtpsessionController(service *goa.Service, ctrl GtpsessionController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateGtpsessionContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateGtpsessionPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/gtpsessions", ctrl.MuxHandler("create", h, unmarshalCreateGtpsessionPayload))
	service.LogInfo("mount", "ctrl", "Gtpsession", "action", "Create", "route", "POST /gtpsessions")
}

// unmarshalCreateGtpsessionPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateGtpsessionPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createGtpsessionPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	payload.Finalize()
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}
