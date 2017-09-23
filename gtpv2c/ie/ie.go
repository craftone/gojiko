package ie

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type ieTypeNum byte

const (
	imsiNum           ieTypeNum = 1
	causeNum          ieTypeNum = 2
	recoveryNum       ieTypeNum = 3
	apnNum            ieTypeNum = 71
	ambrNum           ieTypeNum = 72
	ebiNum            ieTypeNum = 73
	meiNum            ieTypeNum = 75
	msisdnNum         ieTypeNum = 76
	indicationNum     ieTypeNum = 77
	paaNum            ieTypeNum = 79
	ratTypeNum        ieTypeNum = 82
	servingNetworkNum ieTypeNum = 83
	uliNum            ieTypeNum = 86
	fteidNum          ieTypeNum = 87
	pdnTypeNum        ieTypeNum = 99
	apnRestrictionNum ieTypeNum = 127
	selectionModeNum  ieTypeNum = 128
)

type header struct {
	typeNum  ieTypeNum
	length   uint16
	instance byte
}

func newHeader(typeNum ieTypeNum, length uint16, instance byte) (header, error) {
	if instance > 0xf {
		return header{}, fmt.Errorf("instance must be a 4bit number")
	}
	return header{typeNum, length, instance}, nil
}

type IE interface {
	Marshal() []byte
}

func (h *header) marshal(body []byte) []byte {
	bodyLen := uint16(len(body))
	res := make([]byte, 4+bodyLen)
	// make header
	res[0] = byte(h.typeNum)
	binary.BigEndian.PutUint16(res[1:3], bodyLen)
	res[3] = byte(h.instance)
	// copy body
	copy(res[4:], body)
	return res
}

func Unmarshal(buf []byte) (IE, []byte, error) {
	if len(buf) < 4 {
		return nil, buf, errors.New("It needs at least 4 bytes")
	}
	h := header{}
	h.typeNum = ieTypeNum(buf[0])
	h.length = binary.BigEndian.Uint16(buf[1:3])
	msgSize := int(h.length)
	if len(buf) < msgSize {
		return nil, buf, fmt.Errorf("It is too short for the length : %d", h.length)
	}
	h.instance = buf[3] & 0xf

	var msg IE
	var err error
	body := buf[4 : 4+msgSize]

	switch h.typeNum {
	case imsiNum:
		msg, err = unmarshalImsi(h, body)
	case causeNum:
		msg, err = unmarshalCause(h, body)
	case recoveryNum:
		msg, err = unmarshalRecovery(h, body)
	case apnNum:
		msg, err = unmarshalApn(h, body)
	case ambrNum:
		msg, err = unmarshalAmbr(h, body)
	case ebiNum:
		msg, err = unmarshalEbi(h, body)
	case meiNum:
		msg, err = unmarshalMei(h, body)
	case msisdnNum:
		msg, err = unmarshalMsisdn(h, body)
	case indicationNum:
		msg, err = unmarshalIndication(h, body)
	case paaNum:
		msg, err = unmarshalPaa(h, body)
	case ratTypeNum:
		msg, err = unmarshalRatType(h, body)
	case servingNetworkNum:
		msg, err = unmarshalServingNetwork(h, body)
	case uliNum:
		msg, err = unmarshalUli(h, body)
	case fteidNum:
		msg, err = unmarshalFteid(h, body)
	case pdnTypeNum:
		msg, err = unmarshalPdnType(h, body)
	case apnRestrictionNum:
		msg, err = unmarshalApnRestriction(h, body)
	case selectionModeNum:
		msg, err = unmarshalSelectionMode(h, body)
	default:
		return nil, buf, fmt.Errorf("Unknown message type : %d", h.typeNum)
	}
	if err != nil {
		return nil, buf, err
	}
	return msg, buf[4+msgSize:], nil

}
