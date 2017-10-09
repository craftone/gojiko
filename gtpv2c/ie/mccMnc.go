package ie

import (
	"fmt"
)

type mccMnc struct {
	mcc  string
	mnc  string
	tbcd [3]byte
}

func newMccMnc(mcc, mnc string) (mccMnc, error) {
	if len(mcc) != 3 || len(mnc) < 2 || len(mnc) > 3 {
		return mccMnc{}, fmt.Errorf("Invalid length")
	}
	for _, r := range mcc {
		if r < '0' || '9' < r {
			return mccMnc{}, fmt.Errorf("MCC must be a 3-digits number")
		}
	}
	for _, r := range mnc {
		if r < '0' || '9' < r {
			return mccMnc{}, fmt.Errorf("MNC must be a 2- or 3-digits number")
		}
	}

	tbcd := [3]byte{
		(mcc[1]-'0')<<4 + mcc[0] - '0',
		mcc[2] - '0',
		(mnc[1]-'0')<<4 + mnc[0] - '0',
	}
	if len(mnc) == 2 {
		tbcd[1] += 0xf0
	} else {
		tbcd[1] += mnc[2] << 4
	}

	return mccMnc{
		mcc:  mcc,
		mnc:  mnc,
		tbcd: tbcd,
	}, nil
}

func (m mccMnc) copyTo(s []byte) int {
	s[0] = m.tbcd[0]
	s[1] = m.tbcd[1]
	s[2] = m.tbcd[2]
	return 3
}

func unmarshalMccMnc(buf []byte) (mccMnc, []byte, error) {
	if len(buf) < 3 {
		return mccMnc{}, buf, fmt.Errorf("Invalid binary")
	}
	mcc1 := rune(buf[0]&0xf) + '0'
	mcc2 := rune(buf[0]>>4) + '0'
	mcc3 := rune(buf[1]&0xf) + '0'
	mnc1 := rune(buf[2]&0xf) + '0'
	mnc2 := rune(buf[2]>>4) + '0'
	mnc3 := buf[1] >> 4
	mcc := string([]rune{mcc1, mcc2, mcc3})
	mnc := ""
	if mnc3 != 0xf {
		mnc = string([]rune{mnc1, mnc2, rune(mnc3 + '0')})
	} else {
		mnc = string([]rune{mnc1, mnc2})
	}
	mccMncStr, err := newMccMnc(mcc, mnc)
	if err != nil {
		return mccMnc{}, buf, fmt.Errorf("Invalid binary")
	}
	return mccMncStr, buf[3:], nil
}

func (m mccMnc) Mcc() string {
	return m.mcc
}

func (m mccMnc) Mnc() string {
	return m.mnc
}
