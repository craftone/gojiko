package pco

import (
	"encoding/binary"
)

type header struct {
	typeNum pcoTypeNum
	length  byte
}

func newHeader(typeNum pcoTypeNum, length byte) header {
	return header{typeNum, length}
}

func (h *header) marshal(body []byte) []byte {
	bodyLen := byte(len(body))
	res := make([]byte, 3+bodyLen)
	// make header
	binary.BigEndian.PutUint16(res[0:2], uint16(h.typeNum))
	res[2] = bodyLen
	// copy body
	copy(res[3:], body)
	return res
}
