package gtpSessionCmd

type Cmd interface {
	GscType() string
	String() string
}
