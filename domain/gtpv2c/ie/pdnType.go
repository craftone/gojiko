package ie

import "errors"
import "fmt"

type PdnType struct {
	header
	value PdnTypeValue
}

type PdnTypeValue byte

const (
	PdnTypeIPv4   PdnTypeValue = 1
	PdnTypeIPv6   PdnTypeValue = 2
	PdnTypeIPv4v6 PdnTypeValue = 3
)

func (pt PdnTypeValue) isValid() bool {
	if pt == PdnTypeIPv4 || pt == PdnTypeIPv6 || pt == PdnTypeIPv4v6 {
		return true
	}
	return false
}

func NewPdnType(instance byte, value PdnTypeValue) (*PdnType, error) {
	if !value.isValid() {
		return nil, fmt.Errorf("Invalid PDN Type")
	}
	header, err := newHeader(pdnTypeNum, 1, instance)
	if err != nil {
		return nil, err
	}
	return &PdnType{
		header,
		value,
	}, nil
}

func (p *PdnType) Marshal() []byte {
	body := []byte{byte(p.value)}
	return p.header.marshal(body)
}

func unmarshalPdnType(h header, buf []byte) (*PdnType, error) {
	if h.typeNum != pdnTypeNum {
		log.Panic("Invalid type")
	}

	if len(buf) != 1 {
		return nil, errors.New("Invalid binary")
	}

	pdnType := PdnTypeValue(buf[0] & 0x7)
	if !pdnType.isValid() {
		return nil, errors.New("Invalid PDN Type")
	}
	pt, err := NewPdnType(h.instance, pdnType)
	if err != nil {
		return nil, err
	}
	return pt, nil
}

func (p PdnType) Value() PdnTypeValue {
	return p.value
}

func (p PdnType) String() string {
	switch int(p.value) {
	case 1:
		return "IPv4"
	case 2:
		return "IPv6"
	case 3:
		return "IPv4v6"
	}
	return "Unkown"
}
