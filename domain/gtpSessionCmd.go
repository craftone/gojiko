package domain

type gtpSessionCmd interface {
	gscType() string
	String() string
}

type GscResCode int

const (
	GscResOK GscResCode = iota
	GscResNG
	GscResRetryableNG
	GscResTimeout
)

// GscRes respresents the command message that
// replying from the session routine.
type GscRes struct {
	Code    GscResCode
	Msg     string
	Session *GtpSession
}
