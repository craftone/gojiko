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
	StartTime
	EndTime
)

const (
	KTUint64 KeyType = iota
	KTTime
)

var keyTypeMap = map[Key]KeyType{
	SendPackets:        KTUint64,
	SendPacketsSkipped: KTUint64,
	SendBytes:          KTUint64,
	SendBytesSkipped:   KTUint64,
	RecvPackets:        KTUint64,
	RecvPacketsInvalid: KTUint64,
	RecvBytes:          KTUint64,
	RecvBytesInvalid:   KTUint64,
	StartTime:          KTTime,
	EndTime:            KTTime,
}

func Init() {
	log = applog.NewLogger("domain/stats")
	log.Debug("Initialize domain/stats package")
}
