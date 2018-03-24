package domain

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
	Code GsResCode
	Msg  string
	err  error
}
