package gtpv2c

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/craftone/gojiko/gtp"
)

type messageTypeNum byte

const (
	echoRequestNum                   messageTypeNum = 1
	echoResponseNum                  messageTypeNum = 2
	versionNotSupportedIndicationNum messageTypeNum = 3
	createSessionRequestNum          messageTypeNum = 32
	createSessionResponseNum         messageTypeNum = 33
	deleteSessionRequestNum          messageTypeNum = 36
	deleteSessionResponseNum         messageTypeNum = 37
	deleteBearerRequestNum           messageTypeNum = 99
	deleteBearerResponseNum          messageTypeNum = 100
)

type GtpV2cMsg interface {
	Marshal() []byte
}

type header struct {
	version          byte
	messageType      messageTypeNum
	piggybackingFlag bool
	teidFlag         bool
	length           uint16
	teid             gtp.Teid
	seqNum           uint32
}

func newHeader(messageType messageTypeNum, piggybakingFlag, teidFlag bool, teid gtp.Teid, seqNum uint32) header {
	if seqNum > 0xffffff {
		log.Fatal("GTPv2-C's sequence number must be unit24")
	}
	return header{
		version:          2,
		messageType:      messageType,
		piggybackingFlag: piggybakingFlag,
		teidFlag:         teidFlag,
		teid:             teid,
		seqNum:           seqNum,
	}
}

func (h *header) marshal(body []byte) []byte {
	msgLen := uint16(len(body)) + 4
	if h.teidFlag {
		msgLen += 4
	}
	res := make([]byte, msgLen+4)

	// make header
	addr := 0
	res[addr] = byte(h.version << 5)
	if h.piggybackingFlag {
		res[addr] |= (1 << 4)
	}
	if h.teidFlag {
		res[addr] |= (1 << 3)
	}
	addr++

	// Message Type
	res[addr] = byte(h.messageType)
	addr++

	// Message Length
	binary.BigEndian.PutUint16(res[addr:addr+2], msgLen)
	addr += 2

	// TEID
	if h.teidFlag {
		binary.BigEndian.PutUint32(res[addr:addr+4], uint32(h.teid))
		addr += 4
	}

	// Sequence Number
	seqBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(seqBuf, h.seqNum)
	copy(res[addr:], seqBuf[1:4])
	addr += 4

	// copy body
	copy(res[addr:], body)
	return res
}

func Unmarshal(buf []byte) (GtpV2cMsg, []byte, error) {
	if len(buf) < 8 {
		return nil, buf, fmt.Errorf("It needs at least 8 bytes : %v", buf)
	}
	h := header{}
	h.version = buf[0] >> 5
	if h.version != 2 {
		return nil, buf, fmt.Errorf("Version must be 2, but the version is %d", h.version)
	}
	h.piggybackingFlag = (buf[0]&0x10 == 1)
	h.teidFlag = (buf[0]&0x40 == 1)
	h.messageType = messageTypeNum(buf[1])
	h.length = binary.BigEndian.Uint16(buf[2:4])
	msgSize := int(h.length) + 4
	if len(buf) < msgSize {
		return nil, buf, fmt.Errorf("It is too short for the length : %d", h.length)
	}
	idx := 4
	if h.teidFlag {
		if len(buf) < 12 {
			return nil, buf, errors.New("It needs at least 12 bytes")
		}
		h.teid = gtp.Teid(binary.BigEndian.Uint32(buf[idx : idx+4]))
		idx += 4
	}
	h.seqNum = uint32(buf[idx])<<16 + uint32(buf[idx+1])<<8 + uint32(buf[idx+2])
	idx += 4

	var msg GtpV2cMsg
	var err error
	body := buf[idx:msgSize]
	tail := buf[msgSize:]

	switch h.messageType {
	case echoRequestNum:
		msg, err = unmarshalEchoRequest(h, body)
	case echoResponseNum:
		msg, err = unmarshalEchoResponse(h, body)
	default:
		return nil, buf, fmt.Errorf("Unkown message type : %d", h.messageType)
	}
	if err != nil {
		return nil, buf, err
	}
	return msg, tail, nil
}
