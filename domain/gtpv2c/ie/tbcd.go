package ie

import (
	"errors"
	"fmt"
)

type tbcd []byte

func parseTBCD(s string) (tbcd, error) {
	res := make(tbcd, (len(s)+1)/2)
	for i, c := range s {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("Invalid string : %s", s)
		}
		val := byte(c - '0')
		if (i % 2) == 0 {
			res[i/2] = val
		} else {
			res[i/2] += (val << 4)
		}
	}
	if (len(s) % 2) == 1 {
		res[len(s)/2] += 0xf0
	}
	return res, nil
}

func (t tbcd) String() string {
	res := make([]byte, 0, len(t)*2)
	for _, b := range t {
		val1 := b & 0xf
		res = append(res, val1+byte('0'))
		val2 := b >> 4
		if val2 == 0xf {
			break
		}
		res = append(res, val2+byte('0'))
	}
	return string(res)
}

func unmarshalTbcd(buf []byte) (string, error) {
	res := make([]byte, 0, len(buf)*2)
	for _, b := range buf {
		val1 := b & 0xf
		if val1 > 9 {
			return "", errors.New("Invalid tbcd binary")
		}
		res = append(res, val1+byte('0'))
		val2 := b >> 4
		if val2 == 0xf {
			break
		} else if val2 > 9 {
			return "", errors.New("Invalid tbcd binary")
		}
		res = append(res, val2+byte('0'))
	}
	return string(res), nil
}
