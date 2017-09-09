package ie

import (
	"errors"
	"fmt"
)

// Apn struct expresses APN IE
type Apn struct {
	stringIE
}

// ApnMaxLen is maxmum number of digits of APN
const ApnMaxLen = 100

// IsValidAPN determines whether the specific string is valid as an APN.
// All APN must match the following regexp : ^[0-9a-zA-Z-][0-9a-zA-Z-\.]+[0-0a-zA-Z-]$
func IsValidAPN(s string) bool {
	if len(s) < 1 || len(s) > ApnMaxLen {
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
	sie, err := newStringIE(apnNum, uint16(4+len(value)), instance, value, 1, ApnMaxLen)
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
