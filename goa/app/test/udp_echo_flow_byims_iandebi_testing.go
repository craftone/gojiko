// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "gojiko api": udpEchoFlowByIMSIandEBI TestHelpers
//
// Command:
// $ goagen
// --design=github.com/craftone/gojiko/goa/design
// --out=$(GOPATH)/src/github.com/craftone/gojiko/goa
// --regen=true
// --version=v1.3.0

package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/craftone/gojiko/goa/app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/goatest"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// CreateUDPEchoFlowByIMSIandEBIInternalServerError runs the method Create of the given controller with the given parameters and payload.
// It returns the response writer so it's possible to inspect the response headers and the media type struct written to the response.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func CreateUDPEchoFlowByIMSIandEBIInternalServerError(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.UDPEchoFlowByIMSIandEBIController, sgwAddr string, imsi string, ebi int, payload *app.UDPEchoFlowPayload) (http.ResponseWriter, error) {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Validate payload
	err := payload.Validate()
	if err != nil {
		e, ok := err.(goa.ServiceError)
		if !ok {
			panic(err) // bug
		}
		return nil, e
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/sgw/%v/gtpsessions/imsi/%v/ebi/%v/udp_echo_flow", sgwAddr, imsi, ebi),
	}
	req, _err := http.NewRequest("POST", u.String(), nil)
	if _err != nil {
		panic("invalid test " + _err.Error()) // bug
	}
	prms := url.Values{}
	prms["sgwAddr"] = []string{fmt.Sprintf("%v", sgwAddr)}
	prms["imsi"] = []string{fmt.Sprintf("%v", imsi)}
	prms["ebi"] = []string{fmt.Sprintf("%v", ebi)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "UDPEchoFlowByIMSIandEBITest"), rw, req, prms)
	createCtx, __err := app.NewCreateUDPEchoFlowByIMSIandEBIContext(goaCtx, req, service)
	if __err != nil {
		_e, _ok := __err.(goa.ServiceError)
		if !_ok {
			panic("invalid test data " + __err.Error()) // bug
		}
		return nil, _e
	}
	createCtx.Payload = payload

	// Perform action
	__err = ctrl.Create(createCtx)

	// Validate response
	if __err != nil {
		t.Fatalf("controller returned %+v, logs:\n%s", __err, logBuf.String())
	}
	if rw.Code != 500 {
		t.Errorf("invalid response status code: got %+v, expected 500", rw.Code)
	}
	var mt error
	if resp != nil {
		var __ok bool
		mt, __ok = resp.(error)
		if !__ok {
			t.Fatalf("invalid response media: got variable of type %T, value %+v, expected instance of error", resp, resp)
		}
	}

	// Return results
	return rw, mt
}

// CreateUDPEchoFlowByIMSIandEBINotFound runs the method Create of the given controller with the given parameters and payload.
// It returns the response writer so it's possible to inspect the response headers and the media type struct written to the response.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func CreateUDPEchoFlowByIMSIandEBINotFound(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.UDPEchoFlowByIMSIandEBIController, sgwAddr string, imsi string, ebi int, payload *app.UDPEchoFlowPayload) (http.ResponseWriter, error) {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Validate payload
	err := payload.Validate()
	if err != nil {
		e, ok := err.(goa.ServiceError)
		if !ok {
			panic(err) // bug
		}
		return nil, e
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/sgw/%v/gtpsessions/imsi/%v/ebi/%v/udp_echo_flow", sgwAddr, imsi, ebi),
	}
	req, _err := http.NewRequest("POST", u.String(), nil)
	if _err != nil {
		panic("invalid test " + _err.Error()) // bug
	}
	prms := url.Values{}
	prms["sgwAddr"] = []string{fmt.Sprintf("%v", sgwAddr)}
	prms["imsi"] = []string{fmt.Sprintf("%v", imsi)}
	prms["ebi"] = []string{fmt.Sprintf("%v", ebi)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "UDPEchoFlowByIMSIandEBITest"), rw, req, prms)
	createCtx, __err := app.NewCreateUDPEchoFlowByIMSIandEBIContext(goaCtx, req, service)
	if __err != nil {
		_e, _ok := __err.(goa.ServiceError)
		if !_ok {
			panic("invalid test data " + __err.Error()) // bug
		}
		return nil, _e
	}
	createCtx.Payload = payload

	// Perform action
	__err = ctrl.Create(createCtx)

	// Validate response
	if __err != nil {
		t.Fatalf("controller returned %+v, logs:\n%s", __err, logBuf.String())
	}
	if rw.Code != 404 {
		t.Errorf("invalid response status code: got %+v, expected 404", rw.Code)
	}
	var mt error
	if resp != nil {
		var __ok bool
		mt, __ok = resp.(error)
		if !__ok {
			t.Fatalf("invalid response media: got variable of type %T, value %+v, expected instance of error", resp, resp)
		}
	}

	// Return results
	return rw, mt
}

// CreateUDPEchoFlowByIMSIandEBIOK runs the method Create of the given controller with the given parameters and payload.
// It returns the response writer so it's possible to inspect the response headers and the media type struct written to the response.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func CreateUDPEchoFlowByIMSIandEBIOK(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.UDPEchoFlowByIMSIandEBIController, sgwAddr string, imsi string, ebi int, payload *app.UDPEchoFlowPayload) (http.ResponseWriter, *app.Udpechoflow) {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Validate payload
	err := payload.Validate()
	if err != nil {
		e, ok := err.(goa.ServiceError)
		if !ok {
			panic(err) // bug
		}
		t.Errorf("unexpected payload validation error: %+v", e)
		return nil, nil
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/sgw/%v/gtpsessions/imsi/%v/ebi/%v/udp_echo_flow", sgwAddr, imsi, ebi),
	}
	req, _err := http.NewRequest("POST", u.String(), nil)
	if _err != nil {
		panic("invalid test " + _err.Error()) // bug
	}
	prms := url.Values{}
	prms["sgwAddr"] = []string{fmt.Sprintf("%v", sgwAddr)}
	prms["imsi"] = []string{fmt.Sprintf("%v", imsi)}
	prms["ebi"] = []string{fmt.Sprintf("%v", ebi)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "UDPEchoFlowByIMSIandEBITest"), rw, req, prms)
	createCtx, __err := app.NewCreateUDPEchoFlowByIMSIandEBIContext(goaCtx, req, service)
	if __err != nil {
		_e, _ok := __err.(goa.ServiceError)
		if !_ok {
			panic("invalid test data " + __err.Error()) // bug
		}
		t.Errorf("unexpected parameter validation error: %+v", _e)
		return nil, nil
	}
	createCtx.Payload = payload

	// Perform action
	__err = ctrl.Create(createCtx)

	// Validate response
	if __err != nil {
		t.Fatalf("controller returned %+v, logs:\n%s", __err, logBuf.String())
	}
	if rw.Code != 200 {
		t.Errorf("invalid response status code: got %+v, expected 200", rw.Code)
	}
	var mt *app.Udpechoflow
	if resp != nil {
		var __ok bool
		mt, __ok = resp.(*app.Udpechoflow)
		if !__ok {
			t.Fatalf("invalid response media: got variable of type %T, value %+v, expected instance of app.Udpechoflow", resp, resp)
		}
		__err = mt.Validate()
		if __err != nil {
			t.Errorf("invalid response media type: %s", __err)
		}
	}

	// Return results
	return rw, mt
}

// DeleteUDPEchoFlowByIMSIandEBIInternalServerError runs the method Delete of the given controller with the given parameters.
// It returns the response writer so it's possible to inspect the response headers and the media type struct written to the response.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func DeleteUDPEchoFlowByIMSIandEBIInternalServerError(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.UDPEchoFlowByIMSIandEBIController, sgwAddr string, imsi string, ebi int) (http.ResponseWriter, error) {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/sgw/%v/gtpsessions/imsi/%v/ebi/%v/udp_echo_flow", sgwAddr, imsi, ebi),
	}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}
	prms := url.Values{}
	prms["sgwAddr"] = []string{fmt.Sprintf("%v", sgwAddr)}
	prms["imsi"] = []string{fmt.Sprintf("%v", imsi)}
	prms["ebi"] = []string{fmt.Sprintf("%v", ebi)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "UDPEchoFlowByIMSIandEBITest"), rw, req, prms)
	deleteCtx, _err := app.NewDeleteUDPEchoFlowByIMSIandEBIContext(goaCtx, req, service)
	if _err != nil {
		e, ok := _err.(goa.ServiceError)
		if !ok {
			panic("invalid test data " + _err.Error()) // bug
		}
		return nil, e
	}

	// Perform action
	_err = ctrl.Delete(deleteCtx)

	// Validate response
	if _err != nil {
		t.Fatalf("controller returned %+v, logs:\n%s", _err, logBuf.String())
	}
	if rw.Code != 500 {
		t.Errorf("invalid response status code: got %+v, expected 500", rw.Code)
	}
	var mt error
	if resp != nil {
		var _ok bool
		mt, _ok = resp.(error)
		if !_ok {
			t.Fatalf("invalid response media: got variable of type %T, value %+v, expected instance of error", resp, resp)
		}
	}

	// Return results
	return rw, mt
}

// DeleteUDPEchoFlowByIMSIandEBINotFound runs the method Delete of the given controller with the given parameters.
// It returns the response writer so it's possible to inspect the response headers and the media type struct written to the response.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func DeleteUDPEchoFlowByIMSIandEBINotFound(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.UDPEchoFlowByIMSIandEBIController, sgwAddr string, imsi string, ebi int) (http.ResponseWriter, error) {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/sgw/%v/gtpsessions/imsi/%v/ebi/%v/udp_echo_flow", sgwAddr, imsi, ebi),
	}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}
	prms := url.Values{}
	prms["sgwAddr"] = []string{fmt.Sprintf("%v", sgwAddr)}
	prms["imsi"] = []string{fmt.Sprintf("%v", imsi)}
	prms["ebi"] = []string{fmt.Sprintf("%v", ebi)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "UDPEchoFlowByIMSIandEBITest"), rw, req, prms)
	deleteCtx, _err := app.NewDeleteUDPEchoFlowByIMSIandEBIContext(goaCtx, req, service)
	if _err != nil {
		e, ok := _err.(goa.ServiceError)
		if !ok {
			panic("invalid test data " + _err.Error()) // bug
		}
		return nil, e
	}

	// Perform action
	_err = ctrl.Delete(deleteCtx)

	// Validate response
	if _err != nil {
		t.Fatalf("controller returned %+v, logs:\n%s", _err, logBuf.String())
	}
	if rw.Code != 404 {
		t.Errorf("invalid response status code: got %+v, expected 404", rw.Code)
	}
	var mt error
	if resp != nil {
		var _ok bool
		mt, _ok = resp.(error)
		if !_ok {
			t.Fatalf("invalid response media: got variable of type %T, value %+v, expected instance of error", resp, resp)
		}
	}

	// Return results
	return rw, mt
}

// DeleteUDPEchoFlowByIMSIandEBIOK runs the method Delete of the given controller with the given parameters.
// It returns the response writer so it's possible to inspect the response headers and the media type struct written to the response.
// If ctx is nil then context.Background() is used.
// If service is nil then a default service is created.
func DeleteUDPEchoFlowByIMSIandEBIOK(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl app.UDPEchoFlowByIMSIandEBIController, sgwAddr string, imsi string, ebi int) (http.ResponseWriter, *app.Udpechoflow) {
	// Setup service
	var (
		logBuf bytes.Buffer
		resp   interface{}

		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
	)
	if service == nil {
		service = goatest.Service(&logBuf, respSetter)
	} else {
		logger := log.New(&logBuf, "", log.Ltime)
		service.WithLogger(goa.NewLogger(logger))
		newEncoder := func(io.Writer) goa.Encoder { return respSetter }
		service.Encoder = goa.NewHTTPEncoder() // Make sure the code ends up using this decoder
		service.Encoder.Register(newEncoder, "*/*")
	}

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/sgw/%v/gtpsessions/imsi/%v/ebi/%v/udp_echo_flow", sgwAddr, imsi, ebi),
	}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}
	prms := url.Values{}
	prms["sgwAddr"] = []string{fmt.Sprintf("%v", sgwAddr)}
	prms["imsi"] = []string{fmt.Sprintf("%v", imsi)}
	prms["ebi"] = []string{fmt.Sprintf("%v", ebi)}
	if ctx == nil {
		ctx = context.Background()
	}
	goaCtx := goa.NewContext(goa.WithAction(ctx, "UDPEchoFlowByIMSIandEBITest"), rw, req, prms)
	deleteCtx, _err := app.NewDeleteUDPEchoFlowByIMSIandEBIContext(goaCtx, req, service)
	if _err != nil {
		e, ok := _err.(goa.ServiceError)
		if !ok {
			panic("invalid test data " + _err.Error()) // bug
		}
		t.Errorf("unexpected parameter validation error: %+v", e)
		return nil, nil
	}

	// Perform action
	_err = ctrl.Delete(deleteCtx)

	// Validate response
	if _err != nil {
		t.Fatalf("controller returned %+v, logs:\n%s", _err, logBuf.String())
	}
	if rw.Code != 200 {
		t.Errorf("invalid response status code: got %+v, expected 200", rw.Code)
	}
	var mt *app.Udpechoflow
	if resp != nil {
		var _ok bool
		mt, _ok = resp.(*app.Udpechoflow)
		if !_ok {
			t.Fatalf("invalid response media: got variable of type %T, value %+v, expected instance of app.Udpechoflow", resp, resp)
		}
		_err = mt.Validate()
		if _err != nil {
			t.Errorf("invalid response media type: %s", _err)
		}
	}

	// Return results
	return rw, mt
}
