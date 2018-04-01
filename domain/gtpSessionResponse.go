package domain

import (
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
)

type GsResCode int

const (
	GsResOK GsResCode = iota
	GsResNG
	GsResRetryableNG
	GsResTimeout
)

// GsRes respresents the command message that
// replying from the session routine.
type GsRes struct {
	Code  GsResCode
	Value ie.CauseValue
	Msg   string
	err   error
}

func (c GsResCode) String() string {
	switch c {
	case GsResOK:
		return "OK"
	case GsResNG:
		return "NG"
	case GsResRetryableNG:
		return "Retryable NG"
	case GsResTimeout:
		return "Timeout"
	default:
		return "Other"
	}
}
