package ie

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/craftone/gojiko/util"
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
	cgi  *Cgi
	sai  *Sai
	rai  *Rai
	tai  *Tai
	ecgi *Ecgi
	lai  *Lai
}

type Cgi struct {
	mccMnc
	lac uint16
	ci  uint16
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
	cgi, err := NewCgi(mccMnc.mcc, mccMnc.mnc, lac, ci)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return cgi, tail[4:], nil
}

func (c *Cgi) Lac() uint16 {
	return c.lac
}

func (c *Cgi) Ci() uint16 {
	return c.ci
}

type Sai struct {
	mccMnc
	lac uint16
	sac uint16
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
	sai, err := NewSai(mccMnc.mcc, mccMnc.mnc, lac, sac)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return sai, tail[4:], nil
}

func (s *Sai) Lac() uint16 {
	return s.lac
}

func (s *Sai) Sac() uint16 {
	return s.sac
}

type Rai struct {
	mccMnc
	lac uint16
	rac uint16
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
	rai, err := NewRai(mccMnc.mcc, mccMnc.mnc, lac, rac)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return rai, tail[4:], nil
}

func (r *Rai) Lac() uint16 {
	return r.lac
}

func (r *Rai) Rac() uint16 {
	return r.rac
}

type Tai struct {
	mccMnc
	tac uint16
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
	tai, err := NewTai(mccMnc.mcc, mccMnc.mnc, tac)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return tai, tail[2:], nil
}

func (t *Tai) Tac() uint16 {
	return t.tac
}

type Ecgi struct {
	mccMnc
	eci uint32 // actually uint28
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
	ecgi, err := NewEcgi(mccMnc.mcc, mccMnc.mnc, eci)
	if err != nil {
		return nil, buf, fmt.Errorf("Invalid binary")
	}
	return ecgi, tail[4:], nil
}

func (e *Ecgi) Eci() uint32 {
	return e.eci
}

type Lai struct {
	mccMnc
	lac uint16
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
	lai, err := NewLai(mccMnc.mcc, mccMnc.mnc, lac)
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
	return &Uli{
		header,
		uliArg.Cgi,
		uliArg.Sai,
		uliArg.Rai,
		uliArg.Tai,
		uliArg.Ecgi,
		uliArg.Lai,
	}, nil
}

func (r *Uli) Marshal() []byte {
	body := make([]byte, r.length)

	offset := 1
	if r.cgi != nil {
		body[0] = util.SetBit(body[0], 0, true)
		offset += r.cgi.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.cgi.lac)
		offset += 2
		binary.BigEndian.PutUint16(body[offset:], r.cgi.ci)
		offset += 2
	}
	if r.sai != nil {
		body[0] = util.SetBit(body[0], 1, true)
		offset += r.sai.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.sai.lac)
		offset += 2
		binary.BigEndian.PutUint16(body[offset:], r.sai.sac)
		offset += 2
	}
	if r.rai != nil {
		body[0] = util.SetBit(body[0], 2, true)
		offset += r.rai.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.rai.lac)
		offset += 2
		binary.BigEndian.PutUint16(body[offset:], r.rai.rac)
		offset += 2
	}
	if r.tai != nil {
		body[0] = util.SetBit(body[0], 3, true)
		offset += r.tai.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.tai.tac)
		offset += 2
	}
	if r.ecgi != nil {
		body[0] = util.SetBit(body[0], 4, true)
		offset += r.ecgi.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint32(body[offset:], r.ecgi.eci)
		offset += 4
	}
	if r.lai != nil {
		body[0] = util.SetBit(body[0], 5, true)
		offset += r.lai.mccMnc.copyTo(body[offset:])
		binary.BigEndian.PutUint16(body[offset:], r.lai.lac)
		offset += 2
	}

	return r.header.marshal(body)
}

func unmarshalUli(h header, buf []byte) (*Uli, error) {
	if h.typeNum != uliNum {
		log.Panic("Invalid type")
	}

	if len(buf) < 1 {
		return nil, errors.New("Invalid binary")
	}

	var err error
	uliArg := UliArg{}

	tail := buf[1:]
	if util.GetBit(buf[0], 0) {
		uliArg.Cgi, tail, err = unmarshalCgi(tail)
		if err != nil {
			return nil, err
		}
	}
	if util.GetBit(buf[0], 1) {
		uliArg.Sai, tail, err = unmarshalSai(tail)
		if err != nil {
			return nil, err
		}
	}
	if util.GetBit(buf[0], 2) {
		uliArg.Rai, tail, err = unmarshalRai(tail)
		if err != nil {
			return nil, err
		}
	}
	if util.GetBit(buf[0], 3) {
		uliArg.Tai, tail, err = unmarshalTai(tail)
		if err != nil {
			return nil, err
		}
	}
	if util.GetBit(buf[0], 4) {
		uliArg.Ecgi, tail, err = unmarshalEcgi(tail)
		if err != nil {
			return nil, err
		}
	}
	if util.GetBit(buf[0], 5) {
		uliArg.Lai, tail, err = unmarshalLai(tail)
		if err != nil {
			return nil, err
		}
	}
	return NewUli(h.instance, uliArg)
}

func (u *Uli) Cgi() *Cgi {
	return u.cgi
}

func (u *Uli) Sai() *Sai {
	return u.sai
}

func (u *Uli) Rai() *Rai {
	return u.rai
}

func (u *Uli) Tai() *Tai {
	return u.tai
}

func (u *Uli) Ecgi() *Ecgi {
	return u.ecgi
}

func (u *Uli) Lai() *Lai {
	return u.lai
}
