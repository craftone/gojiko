// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "gojiko api": Application Contexts
//
// Command:
// $ goagen
// --design=github.com/craftone/gojiko/goa/design
// --out=$(GOPATH)/src/github.com/craftone/gojiko/goa
// --version=v1.3.0

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// CreateGtpsessionContext provides the gtpsession create action context.
type CreateGtpsessionContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Payload *CreateGtpsessionPayload
}

// NewCreateGtpsessionContext parses the incoming request URL and body, performs validations and creates the
// context used by the gtpsession controller create action.
func NewCreateGtpsessionContext(ctx context.Context, r *http.Request, service *goa.Service) (*CreateGtpsessionContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := CreateGtpsessionContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// createGtpsessionPayload is the gtpsession create action payload.
type createGtpsessionPayload struct {
	// Access Point Name
	Apn *string `form:"apn,omitempty" json:"apn,omitempty" xml:"apn,omitempty"`
	// EPS Bearer ID
	Ebi  *int    `form:"ebi,omitempty" json:"ebi,omitempty" xml:"ebi,omitempty"`
	Imsi *string `form:"imsi,omitempty" json:"imsi,omitempty" xml:"imsi,omitempty"`
	// Mobile Country Code
	Mcc *string `form:"mcc,omitempty" json:"mcc,omitempty" xml:"mcc,omitempty"`
	// Mobile Equipment Identifier
	Mei *string `form:"mei,omitempty" json:"mei,omitempty" xml:"mei,omitempty"`
	// Mobile Network Code
	Mnc    *string `form:"mnc,omitempty" json:"mnc,omitempty" xml:"mnc,omitempty"`
	Msisdn *string `form:"msisdn,omitempty" json:"msisdn,omitempty" xml:"msisdn,omitempty"`
	// SGW GTPv2-C loopback address
	SgwAddr *string `form:"sgwAddr,omitempty" json:"sgwAddr,omitempty" xml:"sgwAddr,omitempty"`
}

// Finalize sets the default values defined in the design.
func (payload *createGtpsessionPayload) Finalize() {
	var defaultEbi = 5
	if payload.Ebi == nil {
		payload.Ebi = &defaultEbi
	}
	var defaultMcc = "440"
	if payload.Mcc == nil {
		payload.Mcc = &defaultMcc
	}
	var defaultMnc = "10"
	if payload.Mnc == nil {
		payload.Mnc = &defaultMnc
	}
}

// Validate runs the validation rules defined in the design.
func (payload *createGtpsessionPayload) Validate() (err error) {
	if payload.SgwAddr == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "sgwAddr"))
	}
	if payload.Apn == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "apn"))
	}
	if payload.Mcc == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "mcc"))
	}
	if payload.Mnc == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "mnc"))
	}
	if payload.Msisdn == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "msisdn"))
	}
	if payload.Mei == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "mei"))
	}
	if payload.Imsi == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "imsi"))
	}
	if payload.Ebi == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "ebi"))
	}
	if payload.Apn != nil {
		if err2 := goa.ValidateFormat(goa.FormatHostname, *payload.Apn); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`raw.apn`, *payload.Apn, goa.FormatHostname, err2))
		}
	}
	if payload.Ebi != nil {
		if *payload.Ebi < 5 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`raw.ebi`, *payload.Ebi, 5, true))
		}
	}
	if payload.Ebi != nil {
		if *payload.Ebi > 15 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`raw.ebi`, *payload.Ebi, 15, false))
		}
	}
	if payload.Imsi != nil {
		if ok := goa.ValidatePattern(`^[0-9]{14,15}$`, *payload.Imsi); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.imsi`, *payload.Imsi, `^[0-9]{14,15}$`))
		}
	}
	if payload.Mcc != nil {
		if ok := goa.ValidatePattern(`^[0-9]{3}$`, *payload.Mcc); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.mcc`, *payload.Mcc, `^[0-9]{3}$`))
		}
	}
	if payload.Mei != nil {
		if ok := goa.ValidatePattern(`^[0-9]{15,16}$`, *payload.Mei); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.mei`, *payload.Mei, `^[0-9]{15,16}$`))
		}
	}
	if payload.Mnc != nil {
		if ok := goa.ValidatePattern(`^[0-9]{2,3}$`, *payload.Mnc); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.mnc`, *payload.Mnc, `^[0-9]{2,3}$`))
		}
	}
	if payload.Msisdn != nil {
		if ok := goa.ValidatePattern(`^[0-9]{12,15}$`, *payload.Msisdn); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.msisdn`, *payload.Msisdn, `^[0-9]{12,15}$`))
		}
	}
	if payload.SgwAddr != nil {
		if err2 := goa.ValidateFormat(goa.FormatIPv4, *payload.SgwAddr); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`raw.sgwAddr`, *payload.SgwAddr, goa.FormatIPv4, err2))
		}
	}
	return
}

// Publicize creates CreateGtpsessionPayload from createGtpsessionPayload
func (payload *createGtpsessionPayload) Publicize() *CreateGtpsessionPayload {
	var pub CreateGtpsessionPayload
	if payload.Apn != nil {
		pub.Apn = *payload.Apn
	}
	if payload.Ebi != nil {
		pub.Ebi = *payload.Ebi
	}
	if payload.Imsi != nil {
		pub.Imsi = *payload.Imsi
	}
	if payload.Mcc != nil {
		pub.Mcc = *payload.Mcc
	}
	if payload.Mei != nil {
		pub.Mei = *payload.Mei
	}
	if payload.Mnc != nil {
		pub.Mnc = *payload.Mnc
	}
	if payload.Msisdn != nil {
		pub.Msisdn = *payload.Msisdn
	}
	if payload.SgwAddr != nil {
		pub.SgwAddr = *payload.SgwAddr
	}
	return &pub
}

// CreateGtpsessionPayload is the gtpsession create action payload.
type CreateGtpsessionPayload struct {
	// Access Point Name
	Apn string `form:"apn" json:"apn" xml:"apn"`
	// EPS Bearer ID
	Ebi  int    `form:"ebi" json:"ebi" xml:"ebi"`
	Imsi string `form:"imsi" json:"imsi" xml:"imsi"`
	// Mobile Country Code
	Mcc string `form:"mcc" json:"mcc" xml:"mcc"`
	// Mobile Equipment Identifier
	Mei string `form:"mei" json:"mei" xml:"mei"`
	// Mobile Network Code
	Mnc    string `form:"mnc" json:"mnc" xml:"mnc"`
	Msisdn string `form:"msisdn" json:"msisdn" xml:"msisdn"`
	// SGW GTPv2-C loopback address
	SgwAddr string `form:"sgwAddr" json:"sgwAddr" xml:"sgwAddr"`
}

// Validate runs the validation rules defined in the design.
func (payload *CreateGtpsessionPayload) Validate() (err error) {
	if payload.SgwAddr == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "sgwAddr"))
	}
	if payload.Apn == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "apn"))
	}
	if payload.Mcc == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "mcc"))
	}
	if payload.Mnc == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "mnc"))
	}
	if payload.Msisdn == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "msisdn"))
	}
	if payload.Mei == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "mei"))
	}
	if payload.Imsi == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "imsi"))
	}

	if err2 := goa.ValidateFormat(goa.FormatHostname, payload.Apn); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`raw.apn`, payload.Apn, goa.FormatHostname, err2))
	}
	if payload.Ebi < 5 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`raw.ebi`, payload.Ebi, 5, true))
	}
	if payload.Ebi > 15 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`raw.ebi`, payload.Ebi, 15, false))
	}
	if ok := goa.ValidatePattern(`^[0-9]{14,15}$`, payload.Imsi); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.imsi`, payload.Imsi, `^[0-9]{14,15}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{3}$`, payload.Mcc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.mcc`, payload.Mcc, `^[0-9]{3}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{15,16}$`, payload.Mei); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.mei`, payload.Mei, `^[0-9]{15,16}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{2,3}$`, payload.Mnc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.mnc`, payload.Mnc, `^[0-9]{2,3}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{12,15}$`, payload.Msisdn); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`raw.msisdn`, payload.Msisdn, `^[0-9]{12,15}$`))
	}
	if err2 := goa.ValidateFormat(goa.FormatIPv4, payload.SgwAddr); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`raw.sgwAddr`, payload.SgwAddr, goa.FormatIPv4, err2))
	}
	return
}

// OK sends a HTTP response with status code 200.
func (ctx *CreateGtpsessionContext) OK(r *Gtpsession) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "application/vnd.gtpsession+json")
	}
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *CreateGtpsessionContext) BadRequest(r error) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	}
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// InternalServerError sends a HTTP response with status code 500.
func (ctx *CreateGtpsessionContext) InternalServerError(r error) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	}
	return ctx.ResponseData.Service.Send(ctx.Context, 500, r)
}
