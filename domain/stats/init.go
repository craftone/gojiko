package stats

import (
	"github.com/craftone/gojiko/applog"
	"github.com/sirupsen/logrus"
)

type Key int
type KeyType int

var log *logrus.Entry

const (
	SendPackets Key = iota
	SendPacketsSkipped
	SendBytes
	SendBytesSkipped
	RecvPackets
	RecvPacketsInvalid
	RecvBytes
	RecvBytesInvalid
)

const (
	KTInt64 KeyType = iota
)

var keyTypeMap = map[Key]KeyType{
	SendPackets:        KTInt64,
	SendPacketsSkipped: KTInt64,
	SendBytes:          KTInt64,
	SendBytesSkipped:   KTInt64,
	RecvPackets:        KTInt64,
	RecvPacketsInvalid: KTInt64,
	RecvBytes:          KTInt64,
	RecvBytesInvalid:   KTInt64,
}

func Init() {
	log = applog.NewLogger("domain/stats")
	log.Debug("Initialize domain/stats package")
}
