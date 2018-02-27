// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "gojiko api": Application Media Types
//
// Command:
// $ goagen
// --design=github.com/craftone/gojiko/goa/design
// --out=$(GOPATH)/src/github.com/craftone/gojiko/goa
// --regen=true
// --version=v1.3.0

package app

import (
	"github.com/goadesign/goa"
)

// A GTP session (default view)
//
// Identifier: application/vnd.gtpsession+json; view=default
type Gtpsession struct {
	// Access Point Name
	Apn string `form:"apn" json:"apn" xml:"apn"`
	// EPS Bearer ID
	Ebi   int               `form:"ebi" json:"ebi" xml:"ebi"`
	Fteid *GtpSessionFTEIDs `form:"fteid" json:"fteid" xml:"fteid"`
	// Session ID in this SGW
	ID   int    `form:"id" json:"id" xml:"id"`
	Imsi string `form:"imsi" json:"imsi" xml:"imsi"`
	// Mobile Country Code
	Mcc string `form:"mcc" json:"mcc" xml:"mcc"`
	// Mobile Equipment Identifier
	Mei string `form:"mei" json:"mei" xml:"mei"`
	// Mobile Network Code
	Mnc    string `form:"mnc" json:"mnc" xml:"mnc"`
	Msisdn string `form:"msisdn" json:"msisdn" xml:"msisdn"`
}

// Validate validates the Gtpsession media type instance.
func (mt *Gtpsession) Validate() (err error) {

	if mt.Apn == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "apn"))
	}
	if mt.Mcc == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "mcc"))
	}
	if mt.Mnc == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "mnc"))
	}
	if mt.Msisdn == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "msisdn"))
	}
	if mt.Mei == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "mei"))
	}
	if mt.Imsi == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "imsi"))
	}

	if mt.Fteid == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "fteid"))
	}
	if err2 := goa.ValidateFormat(goa.FormatHostname, mt.Apn); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`response.apn`, mt.Apn, goa.FormatHostname, err2))
	}
	if mt.Ebi < 5 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`response.ebi`, mt.Ebi, 5, true))
	}
	if mt.Ebi > 15 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`response.ebi`, mt.Ebi, 15, false))
	}
	if mt.Fteid != nil {
		if err2 := mt.Fteid.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if mt.ID < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`response.id`, mt.ID, 0, true))
	}
	if ok := goa.ValidatePattern(`^[0-9]{14,15}$`, mt.Imsi); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.imsi`, mt.Imsi, `^[0-9]{14,15}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{3}$`, mt.Mcc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.mcc`, mt.Mcc, `^[0-9]{3}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{15,16}$`, mt.Mei); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.mei`, mt.Mei, `^[0-9]{15,16}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{2,3}$`, mt.Mnc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.mnc`, mt.Mnc, `^[0-9]{2,3}$`))
	}
	if ok := goa.ValidatePattern(`^[0-9]{12,15}$`, mt.Msisdn); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.msisdn`, mt.Msisdn, `^[0-9]{12,15}$`))
	}
	return
}

// A UDP ECHO flow (default view)
//
// Identifier: application/vnd.udpechoflow+json; view=default
type Udpechoflow struct {
	UDPEchoFlowArg *UDPEchoFlowPayload `form:"UdpEchoFlowArg,omitempty" json:"UdpEchoFlowArg,omitempty" xml:"UdpEchoFlowArg,omitempty"`
}

// Validate validates the Udpechoflow media type instance.
func (mt *Udpechoflow) Validate() (err error) {
	if mt.UDPEchoFlowArg != nil {
		if err2 := mt.UDPEchoFlowArg.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}
