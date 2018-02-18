package ie

import (
	"errors"
	"fmt"
	"strings"
)

// Apn struct expresses APN IE
type Apn struct {
	header
	labels      []string
	totalLength int
}

// ApnMaxLen is maxmum number of digits of APN
const ApnMaxLen = 100

func NewApn(instance byte, str string) (*Apn, error) {
	labels, totalLength := parseAPN(str)
	if len(labels) == 0 {
		return nil, fmt.Errorf("Invalid string for APN : %s", str)
	}

	header, err := newHeader(apnNum, uint16(4+totalLength), instance)
	if err != nil {
		return &Apn{}, err
	}

	return &Apn{
		header:      header,
		labels:      labels,
		totalLength: totalLength,
	}, nil
}

func (i *Apn) Marshal() []byte {
	bytes := make([]byte, 0, i.totalLength)
	for _, label := range i.labels {
		bytes = append(bytes, byte(len(label)))
		bytes = append(bytes, []byte(label)...)
	}
	return i.header.marshal(bytes)
}

func (i *Apn) String() string {
	return strings.Join(i.labels, ".")
}

// parseAPN determines whether the specific string is valid as an APN.
func parseAPN(str string) ([]string, int) {
	totalLength := 0
	labels := strings.Split(str, ".")

	for _, label := range labels {
		totalLength++
		totalLength += len(label)
		if len(label) == 0 || len(label) > 0xFF {
			return []string{}, 0
		}
		// check each charactor
		for _, r := range label {
			// only [0-9a-zA-Z-]
			if !(r == '-') && !('0' <= r && r <= '9') &&
				!('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') {
				return []string{}, 0
			}
		}
		// head and tail should not be '-'.
		if label[0] == '-' || label[len(label)-1] == '-' {
			return []string{}, 0
		}
	}

	if totalLength > ApnMaxLen {
		return []string{}, 0
	}

	return labels, totalLength
}

func unmarshalApn(h header, buf []byte) (*Apn, error) {
	if h.typeNum != apnNum {
		return nil, errors.New("Invalud type")
	}
	labels := make([]string, 0, 6)
	for len(buf) > 0 {
		label, tail := unmarshalLabel(buf)
		if label == "" {
			return nil, errors.New("Invalid bytes")
		}
		labels = append(labels, label)
		buf = tail
	}

	return NewApn(h.instance, strings.Join(labels, "."))
}

func unmarshalLabel(buf []byte) (string, []byte) {
	length := int(buf[0])
	if len(buf) < 1+length {
		return "", []byte{}
	}
	return string(buf[1 : 1+length]), buf[1+length:]
}
