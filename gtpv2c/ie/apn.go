package ie

import (
	"errors"
	"fmt"
)

type Apn struct {
	stringIE
}

const APN_MAX_LEN = 100

func IsValidAPN(s string) bool {
	if len(s) < 1 || len(s) > APN_MAX_LEN {
		return false
	}

	for _, r := range s {
		if !('0' <= r && r <= '9') && !(r == '-') && !(r == '.') &&
			!('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') {
			return false
		}
	}
	if s[0] == '.' || s[len(s)-1] == '.' {
		return false
	}
	return true
}

func NewApn(instance byte, value string) (*Apn, error) {
	if !IsValidAPN(value) {
		return nil, fmt.Errorf("Invalid string for APN : %s", value)
	}
	sie, err := newStringIE(apnNum, 0, instance, value, 1, APN_MAX_LEN)
	if err != nil {
		return nil, err
	}

	return &Apn{
		stringIE: sie,
	}, nil
}

func (i *Apn) Marshal() []byte {
	return i.stringIE.marshal()
}

func unmarshalApn(h header, buf []byte) (*Apn, error) {
	if h.typeNum != apnNum {
		return nil, errors.New("Invalud type")
	}

	return NewApn(h.instance, string(buf))
}
