// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "gojiko api": Application User Types
//
// Command:
// $ goagen
// --design=github.com/craftone/gojiko/goa/design
// --out=$(GOPATH)/src/github.com/craftone/gojiko/goa
// --regen=true
// --version=v1.3.0

package client

import (
	"github.com/goadesign/goa"
)

// fTEID user type.
type fTEID struct {
	IPv4Addr  *string `form:"IPv4 Addr,omitempty" json:"IPv4 Addr,omitempty" xml:"IPv4 Addr,omitempty"`
	IPv4Flag  *bool   `form:"IPv4 Flag,omitempty" json:"IPv4 Flag,omitempty" xml:"IPv4 Flag,omitempty"`
	IPv6Flag  *bool   `form:"IPv6 Flag,omitempty" json:"IPv6 Flag,omitempty" xml:"IPv6 Flag,omitempty"`
	Interface *string `form:"Interface,omitempty" json:"Interface,omitempty" xml:"Interface,omitempty"`
	TEID      *string `form:"TEID,omitempty" json:"TEID,omitempty" xml:"TEID,omitempty"`
}

// Finalize sets the default values for fTEID type instance.
func (ut *fTEID) Finalize() {
	var defaultIPv6Flag = false
	if ut.IPv6Flag == nil {
		ut.IPv6Flag = &defaultIPv6Flag
	}
}

// Validate validates the fTEID type instance.
func (ut *fTEID) Validate() (err error) {
	if ut.IPv4Addr != nil {
		if err2 := goa.ValidateFormat(goa.FormatIPv4, *ut.IPv4Addr); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`request.IPv4 Addr`, *ut.IPv4Addr, goa.FormatIPv4, err2))
		}
	}
	if ut.Interface != nil {
		if !(*ut.Interface == "S5/S8 SGW GTP-U interface" || *ut.Interface == "S5/S8 PGW GTP-U interface" || *ut.Interface == "S5/S8 SGW GTP-C interface" || *ut.Interface == "S5/S8 PGW GTP-C interface") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError(`request.Interface`, *ut.Interface, []interface{}{"S5/S8 SGW GTP-U interface", "S5/S8 PGW GTP-U interface", "S5/S8 SGW GTP-C interface", "S5/S8 PGW GTP-C interface"}))
		}
	}
	if ut.TEID != nil {
		if ok := goa.ValidatePattern(`^0x[0-9A-F]{8}$`, *ut.TEID); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`request.TEID`, *ut.TEID, `^0x[0-9A-F]{8}$`))
		}
	}
	return
}

// Publicize creates FTEID from fTEID
func (ut *fTEID) Publicize() *FTEID {
	var pub FTEID
	if ut.IPv4Addr != nil {
		pub.IPv4Addr = ut.IPv4Addr
	}
	if ut.IPv4Flag != nil {
		pub.IPv4Flag = ut.IPv4Flag
	}
	if ut.IPv6Flag != nil {
		pub.IPv6Flag = *ut.IPv6Flag
	}
	if ut.Interface != nil {
		pub.Interface = ut.Interface
	}
	if ut.TEID != nil {
		pub.TEID = ut.TEID
	}
	return &pub
}

// FTEID user type.
type FTEID struct {
	IPv4Addr  *string `form:"IPv4 Addr,omitempty" json:"IPv4 Addr,omitempty" xml:"IPv4 Addr,omitempty"`
	IPv4Flag  *bool   `form:"IPv4 Flag,omitempty" json:"IPv4 Flag,omitempty" xml:"IPv4 Flag,omitempty"`
	IPv6Flag  bool    `form:"IPv6 Flag" json:"IPv6 Flag" xml:"IPv6 Flag"`
	Interface *string `form:"Interface,omitempty" json:"Interface,omitempty" xml:"Interface,omitempty"`
	TEID      *string `form:"TEID,omitempty" json:"TEID,omitempty" xml:"TEID,omitempty"`
}

// Validate validates the FTEID type instance.
func (ut *FTEID) Validate() (err error) {
	if ut.IPv4Addr != nil {
		if err2 := goa.ValidateFormat(goa.FormatIPv4, *ut.IPv4Addr); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`type.IPv4 Addr`, *ut.IPv4Addr, goa.FormatIPv4, err2))
		}
	}
	if ut.Interface != nil {
		if !(*ut.Interface == "S5/S8 SGW GTP-U interface" || *ut.Interface == "S5/S8 PGW GTP-U interface" || *ut.Interface == "S5/S8 SGW GTP-C interface" || *ut.Interface == "S5/S8 PGW GTP-C interface") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError(`type.Interface`, *ut.Interface, []interface{}{"S5/S8 SGW GTP-U interface", "S5/S8 PGW GTP-U interface", "S5/S8 SGW GTP-C interface", "S5/S8 PGW GTP-C interface"}))
		}
	}
	if ut.TEID != nil {
		if ok := goa.ValidatePattern(`^0x[0-9A-F]{8}$`, *ut.TEID); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`type.TEID`, *ut.TEID, `^0x[0-9A-F]{8}$`))
		}
	}
	return
}

// gtpSessionFTEIDs user type.
type gtpSessionFTEIDs struct {
	PgwCtrlFTEID *fTEID `form:"pgwCtrlFTEID,omitempty" json:"pgwCtrlFTEID,omitempty" xml:"pgwCtrlFTEID,omitempty"`
	PgwDataFTEID *fTEID `form:"pgwDataFTEID,omitempty" json:"pgwDataFTEID,omitempty" xml:"pgwDataFTEID,omitempty"`
	SgwCtrlFTEID *fTEID `form:"sgwCtrlFTEID,omitempty" json:"sgwCtrlFTEID,omitempty" xml:"sgwCtrlFTEID,omitempty"`
	SgwDataFTEID *fTEID `form:"sgwDataFTEID,omitempty" json:"sgwDataFTEID,omitempty" xml:"sgwDataFTEID,omitempty"`
}

// Finalize sets the default values for gtpSessionFTEIDs type instance.
func (ut *gtpSessionFTEIDs) Finalize() {
	if ut.PgwCtrlFTEID != nil {
		var defaultIPv6Flag = false
		if ut.PgwCtrlFTEID.IPv6Flag == nil {
			ut.PgwCtrlFTEID.IPv6Flag = &defaultIPv6Flag
		}
	}
	if ut.PgwDataFTEID != nil {
		var defaultIPv6Flag = false
		if ut.PgwDataFTEID.IPv6Flag == nil {
			ut.PgwDataFTEID.IPv6Flag = &defaultIPv6Flag
		}
	}
	if ut.SgwCtrlFTEID != nil {
		var defaultIPv6Flag = false
		if ut.SgwCtrlFTEID.IPv6Flag == nil {
			ut.SgwCtrlFTEID.IPv6Flag = &defaultIPv6Flag
		}
	}
	if ut.SgwDataFTEID != nil {
		var defaultIPv6Flag = false
		if ut.SgwDataFTEID.IPv6Flag == nil {
			ut.SgwDataFTEID.IPv6Flag = &defaultIPv6Flag
		}
	}
}

// Validate validates the gtpSessionFTEIDs type instance.
func (ut *gtpSessionFTEIDs) Validate() (err error) {
	if ut.PgwCtrlFTEID != nil {
		if err2 := ut.PgwCtrlFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.PgwDataFTEID != nil {
		if err2 := ut.PgwDataFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SgwCtrlFTEID != nil {
		if err2 := ut.SgwCtrlFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SgwDataFTEID != nil {
		if err2 := ut.SgwDataFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// Publicize creates GtpSessionFTEIDs from gtpSessionFTEIDs
func (ut *gtpSessionFTEIDs) Publicize() *GtpSessionFTEIDs {
	var pub GtpSessionFTEIDs
	if ut.PgwCtrlFTEID != nil {
		pub.PgwCtrlFTEID = ut.PgwCtrlFTEID.Publicize()
	}
	if ut.PgwDataFTEID != nil {
		pub.PgwDataFTEID = ut.PgwDataFTEID.Publicize()
	}
	if ut.SgwCtrlFTEID != nil {
		pub.SgwCtrlFTEID = ut.SgwCtrlFTEID.Publicize()
	}
	if ut.SgwDataFTEID != nil {
		pub.SgwDataFTEID = ut.SgwDataFTEID.Publicize()
	}
	return &pub
}

// GtpSessionFTEIDs user type.
type GtpSessionFTEIDs struct {
	PgwCtrlFTEID *FTEID `form:"pgwCtrlFTEID,omitempty" json:"pgwCtrlFTEID,omitempty" xml:"pgwCtrlFTEID,omitempty"`
	PgwDataFTEID *FTEID `form:"pgwDataFTEID,omitempty" json:"pgwDataFTEID,omitempty" xml:"pgwDataFTEID,omitempty"`
	SgwCtrlFTEID *FTEID `form:"sgwCtrlFTEID,omitempty" json:"sgwCtrlFTEID,omitempty" xml:"sgwCtrlFTEID,omitempty"`
	SgwDataFTEID *FTEID `form:"sgwDataFTEID,omitempty" json:"sgwDataFTEID,omitempty" xml:"sgwDataFTEID,omitempty"`
}

// Validate validates the GtpSessionFTEIDs type instance.
func (ut *GtpSessionFTEIDs) Validate() (err error) {
	if ut.PgwCtrlFTEID != nil {
		if err2 := ut.PgwCtrlFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.PgwDataFTEID != nil {
		if err2 := ut.PgwDataFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SgwCtrlFTEID != nil {
		if err2 := ut.SgwCtrlFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if ut.SgwDataFTEID != nil {
		if err2 := ut.SgwDataFTEID.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}
