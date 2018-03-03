package ie

import "errors"
import "fmt"

type SelectionMode struct {
	header
	value byte
}

func NewSelectionMode(instance byte, value byte) (*SelectionMode, error) {
	if value > 3 {
		return nil, fmt.Errorf("Invalid value")
	}
	header, err := newHeader(selectionModeNum, 1, instance)
	if err != nil {
		return nil, err
	}
	return &SelectionMode{
		header,
		value,
	}, nil
}

func (s *SelectionMode) Marshal() []byte {
	body := []byte{s.value}
	return s.header.marshal(body)
}

func unmarshalSelectionMode(h header, buf []byte) (*SelectionMode, error) {
	if h.typeNum != selectionModeNum {
		log.Panic("Invalid type")
	}

	if len(buf) != 1 {
		return nil, errors.New("Invalid binary")
	}

	rec, err := NewSelectionMode(h.instance, buf[0])
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (s *SelectionMode) Value() byte {
	return s.value
}
