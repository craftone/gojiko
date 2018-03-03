package ie

import (
	"errors"

	"github.com/craftone/gojiko/util"
)

type IndicationArg struct {
	DAF   bool
	DTF   bool
	HI    bool
	DFI   bool
	OI    bool
	ISRSI bool
	ISRAI bool
	SGWCI bool
	SQCI  bool
	UIMSI bool
	CFSI  bool
	CRSI  bool
	P     bool
	PT    bool
	SI    bool
	MSV   bool
	ISRAU bool
	CCRSI bool
}

type Indication struct {
	header
	daf   bool
	dtf   bool
	hi    bool
	dfi   bool
	oi    bool
	isrsi bool
	israi bool
	sgwci bool
	sqci  bool
	uimsi bool
	cfsi  bool
	crsi  bool
	p     bool
	pt    bool
	si    bool
	msv   bool
	israu bool
	ccrsi bool
}

func NewIndication(instance byte, indicationArg IndicationArg) (*Indication, error) {
	header, err := newHeader(indicationNum, 3, instance)
	if err != nil {
		return nil, err
	}
	return &Indication{
		header: header,
		daf:    indicationArg.DAF,
		dtf:    indicationArg.DTF,
		hi:     indicationArg.HI,
		dfi:    indicationArg.DFI,
		oi:     indicationArg.OI,
		isrsi:  indicationArg.ISRSI,
		israi:  indicationArg.ISRAI,
		sgwci:  indicationArg.SGWCI,
		sqci:   indicationArg.SQCI,
		uimsi:  indicationArg.UIMSI,
		cfsi:   indicationArg.CFSI,
		crsi:   indicationArg.CRSI,
		p:      indicationArg.P,
		pt:     indicationArg.PT,
		si:     indicationArg.SI,
		msv:    indicationArg.MSV,
		israu:  indicationArg.ISRAU,
		ccrsi:  indicationArg.CCRSI,
	}, nil
}

func (i *Indication) Marshal() []byte {
	// on 3gpp R10 indication size is 3 + alpha
	// on 3gpp R14 indication size is 7 + alpha
	body := make([]byte, 7)
	body[0] = util.SetBit(body[0], 7, i.daf)
	body[0] = util.SetBit(body[0], 6, i.dtf)
	body[0] = util.SetBit(body[0], 5, i.hi)
	body[0] = util.SetBit(body[0], 4, i.dfi)
	body[0] = util.SetBit(body[0], 3, i.oi)
	body[0] = util.SetBit(body[0], 2, i.isrsi)
	body[0] = util.SetBit(body[0], 1, i.israi)
	body[0] = util.SetBit(body[0], 0, i.sgwci)
	body[1] = util.SetBit(body[1], 7, i.sqci)
	body[1] = util.SetBit(body[1], 6, i.uimsi)
	body[1] = util.SetBit(body[1], 5, i.cfsi)
	body[1] = util.SetBit(body[1], 4, i.crsi)
	body[1] = util.SetBit(body[1], 3, i.p)
	body[1] = util.SetBit(body[1], 2, i.pt)
	body[1] = util.SetBit(body[1], 1, i.si)
	body[1] = util.SetBit(body[1], 0, i.msv)
	body[2] = util.SetBit(body[2], 1, i.israu)
	body[2] = util.SetBit(body[2], 0, i.ccrsi)

	return i.header.marshal(body)
}

func unmarshalIndication(h header, buf []byte) (*Indication, error) {
	if h.typeNum != indicationNum {
		log.Panic("Invalid type")
	}

	if len(buf) < 3 {
		return nil, errors.New("Invalid binary")
	}

	indicationArg := IndicationArg{}

	indicationArg.DAF = util.GetBit(buf[0], 7)
	indicationArg.DTF = util.GetBit(buf[0], 6)
	indicationArg.HI = util.GetBit(buf[0], 5)
	indicationArg.DFI = util.GetBit(buf[0], 4)
	indicationArg.OI = util.GetBit(buf[0], 3)
	indicationArg.ISRSI = util.GetBit(buf[0], 2)
	indicationArg.ISRAI = util.GetBit(buf[0], 1)
	indicationArg.SGWCI = util.GetBit(buf[0], 0)
	indicationArg.SQCI = util.GetBit(buf[1], 7)
	indicationArg.UIMSI = util.GetBit(buf[1], 6)
	indicationArg.CFSI = util.GetBit(buf[1], 5)
	indicationArg.CRSI = util.GetBit(buf[1], 4)
	indicationArg.P = util.GetBit(buf[1], 3)
	indicationArg.PT = util.GetBit(buf[1], 2)
	indicationArg.SI = util.GetBit(buf[1], 1)
	indicationArg.MSV = util.GetBit(buf[1], 0)
	indicationArg.ISRAU = util.GetBit(buf[2], 1)
	indicationArg.CCRSI = util.GetBit(buf[2], 0)
	return NewIndication(h.instance, indicationArg)
}

func (i *Indication) DAF() bool {
	return i.daf
}

func (i *Indication) DTF() bool {
	return i.dtf
}

func (i *Indication) HI() bool {
	return i.hi
}

func (i *Indication) DFI() bool {
	return i.dfi
}

func (i *Indication) OI() bool {
	return i.oi
}

func (i *Indication) ISRSI() bool {
	return i.isrsi
}

func (i *Indication) ISRAI() bool {
	return i.israi
}

func (i *Indication) SGWCI() bool {
	return i.sgwci
}

func (i *Indication) SQCI() bool {
	return i.sqci
}

func (i *Indication) UIMSI() bool {
	return i.uimsi
}

func (i *Indication) CFSI() bool {
	return i.cfsi
}

func (i *Indication) CRSI() bool {
	return i.crsi
}

func (i *Indication) P() bool {
	return i.p
}

func (i *Indication) PT() bool {
	return i.pt
}

func (i *Indication) SI() bool {
	return i.si
}

func (i *Indication) MSV() bool {
	return i.msv
}

func (i *Indication) ISRAU() bool {
	return i.israu
}

func (i *Indication) CCRSI() bool {
	return i.ccrsi
}
