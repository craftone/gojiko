package gtpSessionCmd

import (
	"fmt"
)

type CSResCode byte

const (
	CSResOK CSResCode = iota
	CSResNG
	CSResTimeout
)

// CreateSessionRes respresents the command message that
// replying from the session routine.
type CreateSessionRes struct {
	Code CSResCode
	Msg  string
}

func (g CreateSessionRes) GscType() string {
	return "CreateSessionRes"
}

func (g CreateSessionRes) String() string {
	return fmt.Sprintf("TYPE=%s ResCode=%d MSG=%s", g.GscType(), g.Code, g.Msg)
}
