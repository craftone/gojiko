package ie

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

type UliArg struct {
	Cgi  *Cgi
	Sai  *Sai
	Rai  *Rai
	Tai  *Tai
	Ecgi *Ecgi
	Lai  *Lai
}

type Uli struct {
	header
	UliArg
}

type Cgi struct {
	mccMnc
	Lac uint16
	Ci  uint16
}

func NewCgi(mcc, mnc string, lac, ci uint16) (*Cgi, error) {
	mccMnc, err := newMccMnc(mcc, mnc)
	if err != nil {
		return nil, err
	}
	return &Cgi{mccMnc, lac, ci}, nil
}

func unmarshalCgi(buf []byte) (*Cgi, []byte, error) {
	mccMnc, tail, err := unmarshalMccMnc(buf)
	if err != nil {
		return nil, buf, err
	}
	if len(tail) < 4 {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	lac := binary.BigEndian.Uint16(tail[0:2])
	ci := binary.BigEndian.Uint16(tail[2:4])
	cgi, err := NewCgi(mccMnc.Mcc, mccMnc.Mnc, lac, ci)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return cgi, tail[4:], nil
}

type Sai struct {
	mccMnc
	Lac uint16
	Sac uint16
}

func NewSai(mcc, mnc string, lac, sac uint16) (*Sai, error) {
	mccMnc, err := newMccMnc(mcc, mnc)
	if err != nil {
		return nil, err
	}
	return &Sai{mccMnc, lac, sac}, nil
}

func unmarshalSai(buf []byte) (*Sai, []byte, error) {
	mccMnc, tail, err := unmarshalMccMnc(buf)
	if err != nil {
		return nil, buf, err
	}
	if len(tail) < 4 {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	lac := binary.BigEndian.Uint16(tail[0:2])
	sac := binary.BigEndian.Uint16(tail[2:4])
	sai, err := NewSai(mccMnc.Mcc, mccMnc.Mnc, lac, sac)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return sai, tail[4:], nil
}

type Rai struct {
	mccMnc
	Lac uint16
	Rac uint16
}

func NewRai(mcc, mnc string, lac, rac uint16) (*Rai, error) {
	mccMnc, err := newMccMnc(mcc, mnc)
	if err != nil {
		return nil, err
	}
	return &Rai{mccMnc, lac, rac}, nil
}

func unmarshalRai(buf []byte) (*Rai, []byte, error) {
	mccMnc, tail, err := unmarshalMccMnc(buf)
	if err != nil {
		return nil, buf, err
	}
	if len(tail) < 4 {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	lac := binary.BigEndian.Uint16(tail[0:2])
	rac := binary.BigEndian.Uint16(tail[2:4])
	rai, err := NewRai(mccMnc.Mcc, mccMnc.Mnc, lac, rac)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return rai, tail[4:], nil
}

type Tai struct {
	mccMnc
	Tac uint16
}

func NewTai(mcc, mnc string, tac uint16) (*Tai, error) {
	mccMnc, err := newMccMnc(mcc, mnc)
	if err != nil {
		return nil, err
	}
	return &Tai{mccMnc, tac}, nil
}

func unmarshalTai(buf []byte) (*Tai, []byte, error) {
	mccMnc, tail, err := unmarshalMccMnc(buf)
	if err != nil {
		return nil, buf, err
	}
	if len(tail) < 2 {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	tac := binary.BigEndian.Uint16(tail[0:2])
	tai, err := NewTai(mccMnc.Mcc, mccMnc.Mnc, tac)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return tai, tail[2:], nil
}

type Ecgi struct {
	mccMnc
	Eci uint32 // actually uint28
}

func NewEcgi(mcc, mnc string, eci uint32) (*Ecgi, error) {
	mccMnc, err := newMccMnc(mcc, mnc)
	if err != nil {
		return nil, err
	}
	if eci > 0xfffffff {
		return nil, fmt.Errorf("Invalid ECI, too big : %X", eci)
	}
	return &Ecgi{mccMnc, eci}, nil
}

func unmarshalEcgi(buf []byte) (*Ecgi, []byte, error) {
	mccMnc, tail, err := unmarshalMccMnc(buf)
	if err != nil {
		return nil, buf, err
	}
	if len(tail) < 4 {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	eci := binary.BigEndian.Uint32(tail[0:4]) & 0xfffffff
	ecgi, err := NewEcgi(mccMnc.Mcc, mccMnc.Mnc, eci)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return ecgi, tail[4:], nil
}

type Lai struct {
	mccMnc
	Lac uint16
}

func NewLai(mcc, mnc string, lai uint16) (*Lai, error) {
	mccMnc, err := newMccMnc(mcc, mnc)
	if err != nil {
		return nil, err
	}
	return &Lai{mccMnc, lai}, nil
}

func unmarshalLai(buf []byte) (*Lai, []byte, error) {
	mccMnc, tail, err := unmarshalMccMnc(buf)
	if err != nil {
		return nil, buf, err
	}
	if len(tail) < 2 {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	lac := binary.BigEndian.Uint16(tail[0:2])
	lai, err := NewLai(mccMnc.Mcc, mccMnc.Mnc, lac)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return lai, tail[2:], nil
}

func NewUli(instance byte, uliArg UliArg) (*Uli, error) {
	length := 1
	if uliArg.Cgi != nil {
		length += 7
	}
	if uliArg.Sai != nil {
		length += 7
	}
	if uliArg.Rai != nil {
		length += 7
	}
	if uliArg.Tai != nil {
		length += 5
	}
	if uliArg.Ecgi != nil {
		length += 7
	}
	header, err := newHeader(uliNum, uint16(length), instance)
	if err != nil {
		return nil, err
	}
	return &Uli{header, uliArg}, nil
}

func (r *Uli) Marshal() []byte {
	body := make([]byte, r.length)

	offset := 1
	if r.Cgi != nil {
		body[0] = setBit(body[0], 0, true)
		offset += r.Cgi.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.Cgi.Lac)
		offset += 2
		binary.BigEndian.PutUint16(body[offset:], r.Cgi.Ci)
		offset += 2
	}
	if r.Sai != nil {
		body[0] = setBit(body[0], 1, true)
		offset += r.Sai.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.Sai.Lac)
		offset += 2
		binary.BigEndian.PutUint16(body[offset:], r.Sai.Sac)
		offset += 2
	}
	if r.Rai != nil {
		body[0] = setBit(body[0], 2, true)
		offset += r.Rai.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.Rai.Lac)
		offset += 2
		binary.BigEndian.PutUint16(body[offset:], r.Rai.Rac)
		offset += 2
	}
	if r.Tai != nil {
		body[0] = setBit(body[0], 3, true)
		offset += r.Tai.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.Tai.Tac)
		offset += 2
	}
	if r.Ecgi != nil {
		body[0] = setBit(body[0], 4, true)
		offset += r.Ecgi.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint32(body[offset:], r.Ecgi.Eci)
		offset += 4
	}
	if r.Lai != nil {
		body[0] = setBit(body[0], 5, true)
		offset += r.Lai.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.Lai.Lac)
		offset += 2
	}

	return r.header.marshal(body)
}

func unmarshalUli(h header, buf []byte) (*Uli, error) {
	if h.typeNum != uliNum {
		log.Fatal("Invalud type")
	}

	if len(buf) < 1 {
		return nil, errors.New("Invalid binary")
	}

	var err error
	uliArg := UliArg{}

	tail := buf[1:]
	if getBit(buf[0], 0) {
		uliArg.Cgi, tail, err = unmarshalCgi(tail)
		if err != nil {
			return nil, err
		}
	}
	if getBit(buf[0], 1) {
		uliArg.Sai, tail, err = unmarshalSai(tail)
		if err != nil {
			return nil, err
		}
	}
	if getBit(buf[0], 2) {
		uliArg.Rai, tail, err = unmarshalRai(tail)
		if err != nil {
			return nil, err
		}
	}
	if getBit(buf[0], 3) {
		uliArg.Tai, tail, err = unmarshalTai(tail)
		if err != nil {
			return nil, err
		}
	}
	if getBit(buf[0], 4) {
		uliArg.Ecgi, tail, err = unmarshalEcgi(tail)
		if err != nil {
			return nil, err
		}
	}
	if getBit(buf[0], 5) {
		uliArg.Lai, tail, err = unmarshalLai(tail)
		if err != nil {
			return nil, err
		}
	}
	return NewUli(h.instance, uliArg)
}
