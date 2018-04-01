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

// A GTP session (default view)
//
// Identifier: application/vnd.gtpsession+json; view=default
type Gtpsession struct {
	// Access Point Name
	Apn string `form:"apn" json:"apn" xml:"apn"`
	// EPS Bearer ID
	Ebi   int               `form:"ebi" json:"ebi" xml:"ebi"`
	Fteid *GtpSessionFTEIDs `form:"fteid,omitempty" json:"fteid,omitempty" xml:"fteid,omitempty"`
	Imsi  string            `form:"imsi" json:"imsi" xml:"imsi"`
	// Mobile Country Code
	Mcc string `form:"mcc" json:"mcc" xml:"mcc"`
	// Mobile Equipment Identifier
	Mei string `form:"mei" json:"mei" xml:"mei"`
	// Mobile Network Code
	Mnc    string `form:"mnc" json:"mnc" xml:"mnc"`
	Msisdn string `form:"msisdn" json:"msisdn" xml:"msisdn"`
	// PDN Address Allocation
	Paa string `form:"paa" json:"paa" xml:"paa"`
	// Session ID in this SGW
	Sid int `form:"sid" json:"sid" xml:"sid"`
}

// Validate validates the Gtpsession media type instance.
func (mt *Gtpsession) Validate() (err error) {
	if mt.Apn == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "apn"))
	}

	if mt.Imsi == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "imsi"))
	}
	if mt.Mcc == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "mcc"))
	}
	if mt.Mnc == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "mnc"))
	}
	if mt.Mei == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "mei"))
	}
	if mt.Msisdn == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "msisdn"))
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
	if err2 := goa.ValidateFormat(goa.FormatIPv4, mt.Paa); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`response.paa`, mt.Paa, goa.FormatIPv4, err2))
	}
	if mt.Sid < 0 {
		err = goa.MergeErrors(err, goa.InvalidRangeError(`response.sid`, mt.Sid, 0, true))
	}
	return
}

// GTPv2-C Cause (default view)
//
// Identifier: application/vnd.gtpv2c.cause+json; view=default
type Gtpv2cCause struct {
	// Detail of return code from PGW
	Detail string `form:"detail" json:"detail" xml:"detail"`
	// Type of return code from PGW
	Type string `form:"type" json:"type" xml:"type"`
	// GTPv2-C response Cause Value
	Value int `form:"value" json:"value" xml:"value"`
}

// Validate validates the Gtpv2cCause media type instance.
func (mt *Gtpv2cCause) Validate() (err error) {
	if mt.Type == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "type"))
	}

	if mt.Detail == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "detail"))
	}
	return
}

// GTPv2-C Crease Session Response (default view)
//
// Identifier: application/vnd.gtpv2c.csres+json; view=default
type Gtpv2cCsres struct {
	Cause       *Gtpv2cCause `form:"cause" json:"cause" xml:"cause"`
	SessionInfo *Gtpsession  `form:"sessionInfo" json:"sessionInfo" xml:"sessionInfo"`
}

// Validate validates the Gtpv2cCsres media type instance.
func (mt *Gtpv2cCsres) Validate() (err error) {
	if mt.Cause == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "cause"))
	}
	if mt.SessionInfo == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "sessionInfo"))
	}
	if mt.Cause != nil {
		if err2 := mt.Cause.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if mt.SessionInfo != nil {
		if err2 := mt.SessionInfo.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// A UDP ECHO flow (default view)
//
// Identifier: application/vnd.udpechoflow+json; view=default
type Udpechoflow struct {
	Param *UDPEchoFlowPayload `form:"param,omitempty" json:"param,omitempty" xml:"param,omitempty"`
}

// Validate validates the Udpechoflow media type instance.
func (mt *Udpechoflow) Validate() (err error) {
	if mt.Param != nil {
		if err2 := mt.Param.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// A UDP ECHO flow (withStats view)
//
// Identifier: application/vnd.udpechoflow+json; view=withStats
type UdpechoflowWithStats struct {
	Param *UDPEchoFlowPayload `form:"param,omitempty" json:"param,omitempty" xml:"param,omitempty"`
	Stats *SendRecvStatistics `form:"stats,omitempty" json:"stats,omitempty" xml:"stats,omitempty"`
}

// Validate validates the UdpechoflowWithStats media type instance.
func (mt *UdpechoflowWithStats) Validate() (err error) {
	if mt.Param != nil {
		if err2 := mt.Param.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if mt.Stats != nil {
		if err2 := mt.Stats.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}
