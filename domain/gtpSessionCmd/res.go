package gtpSessionCmd

type ResCode int

const (
	ResOK      ResCode = 200
	ResNG      ResCode = 400
	ResTimeout ResCode = 500
)

// Res respresents the command message that
// replying from the session routine.
type Res struct {
	Code ResCode
	Msg  string
}
