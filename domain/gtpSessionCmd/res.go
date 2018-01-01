package gtpSessionCmd

type ResCode int

const (
	ResOK ResCode = iota
	ResNG
	ResBusy
	ResTimeout
)

// Res respresents the command message that
// replying from the session routine.
type Res struct {
	Code ResCode
	Msg  string
}
