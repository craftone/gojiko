// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "gojiko api": Application Media Types
//
// Command:
// $ goagen
// --design=github.com/craftone/gojiko/goa/design
// --out=$(GOPATH)/src/github.com/craftone/gojiko/goa
// --version=v1.3.0

package app

import (
	"github.com/goadesign/goa"
)

// A gtp session (default view)
//
// Identifier: application/vnd.gtpsession+json; view=default
type Gtpsession struct {
	// PGW's Access Point Name
	Apn *string `form:"apn,omitempty" json:"apn,omitempty" xml:"apn,omitempty"`
	// EPS Bearer ID
	Ebi   int               `form:"ebi" json:"ebi" xml:"ebi"`
	Fteid *GtpSessionFTEIDs `form:"fteid,omitempty" json:"fteid,omitempty" xml:"fteid,omitempty"`
	// Session ID in this SGW
	ID   *int    `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	Imsi *string `form:"imsi,omitempty" json:"imsi,omitempty" xml:"imsi,omitempty"`
	// PGW's Mobile Country Code
	Mcc string `form:"mcc" json:"mcc" xml:"mcc"`
	// Mobile Equipment Identifier
	Mei *string `form:"mei,omitempty" json:"mei,omitempty" xml:"mei,omitempty"`
	// PGW's Mobile Network Code
	Mnc    string  `form:"mnc" json:"mnc" xml:"mnc"`
	Msisdn *string `form:"msisdn,omitempty" json:"msisdn,omitempty" xml:"msisdn,omitempty"`
}

// Validate validates the Gtpsession media type instance.
func (mt *Gtpsession) Validate() (err error) {
	if mt.Apn != nil {
		if err2 := goa.ValidateFormat(goa.FormatHostname, *mt.Apn); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`response.apn`, *mt.Apn, goa.FormatHostname, err2))
		}
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
	if mt.ID != nil {
		if *mt.ID < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError(`response.id`, *mt.ID, 0, true))
		}
	}
	if mt.Imsi != nil {
		if ok := goa.ValidatePattern(`^[0-9]{14,15}$`, *mt.Imsi); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`response.imsi`, *mt.Imsi, `^[0-9]{14,15}$`))
		}
	}
	if ok := goa.ValidatePattern(`^[0-9]{3}$`, mt.Mcc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.mcc`, mt.Mcc, `^[0-9]{3}$`))
	}
	if mt.Mei != nil {
		if ok := goa.ValidatePattern(`^[0-9]{15,16}$`, *mt.Mei); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`response.mei`, *mt.Mei, `^[0-9]{15,16}$`))
		}
	}
	if ok := goa.ValidatePattern(`^[0-9]{2,3}$`, mt.Mnc); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.mnc`, mt.Mnc, `^[0-9]{2,3}$`))
	}
	if mt.Msisdn != nil {
		if ok := goa.ValidatePattern(`^[0-9]{12,15}$`, *mt.Msisdn); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`response.msisdn`, *mt.Msisdn, `^[0-9]{12,15}$`))
		}
	}
	return
}
