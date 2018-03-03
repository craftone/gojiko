package gtpv2c

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/util"
)

type MessageTypeNum byte

const (
	EchoRequestNum                   MessageTypeNum = 1
	EchoResponseNum                  MessageTypeNum = 2
	VersionNotSupportedIndicationNum MessageTypeNum = 3
	CreateSessionRequestNum          MessageTypeNum = 32
	CreateSessionResponseNum         MessageTypeNum = 33
	DeleteSessionRequestNum          MessageTypeNum = 36
	DeleteSessionResponseNum         MessageTypeNum = 37
	DeleteBearerRequestNum           MessageTypeNum = 99
	DeleteBearerResponseNum          MessageTypeNum = 100
)

type GtpV2cMsg interface {
	Marshal() []byte
	TeidFlag() bool
	Teid() gtp.Teid
	SeqNum() uint32
}

type header struct {
	version          byte
	messageType      MessageTypeNum
	piggybackingFlag bool
	teidFlag         bool
	length           uint16
	teid             gtp.Teid
	seqNum           uint32
}

func newHeader(messageType MessageTypeNum, piggybakingFlag, teidFlag bool, teid gtp.Teid, seqNum uint32) header {
	if seqNum > 0xffffff {
		log.Panicf("GTPv2-C's sequence number must be unit24 but %d", seqNum)
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
	res[addr] = util.SetBit(res[addr], 4, h.piggybackingFlag)
	res[addr] = util.SetBit(res[addr], 3, h.teidFlag)
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

	// Unmarshal header part
	h := header{}
	h.version = buf[0] >> 5
	if h.version != 2 {
		return nil, buf, fmt.Errorf("Version must be 2, but the version is %d", h.version)
	}
	h.piggybackingFlag = util.GetBit(buf[0], 4)
	h.teidFlag = util.GetBit(buf[0], 3)
	h.messageType = MessageTypeNum(buf[1])
	h.length = binary.BigEndian.Uint16(buf[2:4])
	msgSize := int(h.length) + 4
	if len(buf) < msgSize {
		return nil, buf, fmt.Errorf("The binary size %d is too short for the length : %d", len(buf), msgSize)
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

	// Unmarshal body part
	switch h.messageType {
	case EchoRequestNum:
		msg, err = unmarshalEchoRequest(h, body)
	case EchoResponseNum:
		msg, err = unmarshalEchoResponse(h, body)
	case CreateSessionRequestNum:
		msg, err = unmarshalCreateSessionRequest(h, body)
	case CreateSessionResponseNum:
		msg, err = unmarshalCreateSessionResponse(h, body)
	case DeleteBearerRequestNum:
		msg, err = unmarshalDeleteBearerRequest(h, body)
	case DeleteBearerResponseNum:
		msg, err = unmarshalDeleteBearerResponse(h, body)
	default:
		return nil, buf, fmt.Errorf("Unkown message type : %d", h.messageType)
	}
	if err != nil {
		return nil, buf, err
	}
	return msg, tail, nil
}

//
// Getters
//

func (h *header) TeidFlag() bool {
	return h.teidFlag
}

func (h *header) Teid() gtp.Teid {
	return h.teid
}

func (h *header) SeqNum() uint32 {
	return h.seqNum
}
