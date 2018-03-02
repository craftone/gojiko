package gtpv1u

import (
	"encoding/binary"

	"github.com/craftone/gojiko/gtp"
	"github.com/craftone/gojiko/util"
)

type MessageTypeNum byte

const (
	EchoRequestNum     MessageTypeNum = 1
	EchoResponseNum    MessageTypeNum = 2
	ErrorIndicationNum MessageTypeNum = 26
	GpduNum            MessageTypeNum = 255
)

type header struct {
	version             byte
	protocolType        bool
	extentionHeaderFlag bool
	sequenceNumberFlag  bool
	npduNumberFlag      bool
	messageType         MessageTypeNum
	length              uint16
	teid                gtp.Teid
	seqNum              uint16
}

func newHeader(messageType MessageTypeNum, teid gtp.Teid, seqNum uint16) header {
	return header{
		version:             1,
		protocolType:        true, //true:GTP, false:GTP'
		extentionHeaderFlag: false,
		sequenceNumberFlag:  true, //always true, at this time
		npduNumberFlag:      false,
		messageType:         messageType,
		teid:                teid,
		seqNum:              seqNum,
	}
}

func (h header) marshal(body []byte) []byte {
	res := make([]byte, 12+len(body))
	res[0] = h.version << 4
	res[0] = util.SetBit(res[0], 4, h.protocolType)
	res[0] = util.SetBit(res[0], 2, h.extentionHeaderFlag)
	res[0] = util.SetBit(res[0], 1, h.sequenceNumberFlag)
	res[0] = util.SetBit(res[0], 0, h.npduNumberFlag)
	res[1] = byte(h.messageType)
	binary.BigEndian.PutUint16(res[2:], uint16(4+len(body)))
	binary.BigEndian.PutUint32(res[4:], uint32(h.teid))
	binary.BigEndian.PutUint16(res[8:], h.seqNum)
	if len(body) > 0 {
		copy(res[12:], body)
	}
	return res
}