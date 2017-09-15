package ie

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaa(t *testing.T) {
	paa, _ := NewPaa(0, PdnTypeIPv4, net.IPv4(1, 2, 3, 4), nil)
	assert.Equal(t, paaNum, paa.typeNum)
	assert.Equal(t, PdnTypeIPv4, paa.Value)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), paa.ipv4)
	assert.Nil(t, paa.ipv6)

	paa, _ = NewPaa(0, PdnTypeIPv6, nil, net.ParseIP("2001:db8::68"))
	assert.Equal(t, PdnTypeIPv6, paa.Value)
	assert.Nil(t, paa.ipv4)
	assert.Equal(t, net.ParseIP("2001:db8::68").To16(), paa.ipv6)

	paa, _ = NewPaa(0, PdnTypeIPv4v6, net.IPv4(1, 2, 3, 4), net.ParseIP("2001:db8::68"))
	assert.Equal(t, PdnTypeIPv4v6, paa.Value)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), paa.ipv4)
	assert.Equal(t, net.ParseIP("2001:db8::68").To16(), paa.ipv6)

	// Both IPs are nil
	_, err := NewPaa(0, PdnTypeIPv4, nil, nil)
	assert.Error(t, err)
	_, err = NewPaa(0, PdnTypeIPv6, nil, nil)
	assert.Error(t, err)
	_, err = NewPaa(0, PdnTypeIPv4v6, nil, nil)
	assert.Error(t, err)

	// PdnType mismatch
	_, err = NewPaa(0, PdnTypeIPv4, nil, net.IPv4(1, 2, 3, 4))
	assert.Error(t, err)
	_, err = NewPaa(0, PdnTypeIPv6, net.IPv4(1, 2, 3, 4), nil)
	assert.Error(t, err)
	_, err = NewPaa(0, PdnTypeIPv4v6, nil, net.IPv4(1, 2, 3, 4))
	assert.Error(t, err)
	_, err = NewPaa(0, PdnTypeIPv4v6, net.IPv4(1, 2, 3, 4), nil)
	assert.Error(t, err)
}

func TestPaa_marshal(t *testing.T) {
	paa, _ := NewPaa(0, PdnTypeIPv4, net.IPv4(1, 2, 3, 4), nil)
	paaBin := paa.Marshal()
	assert.Equal(t, []byte{0x4f, 0, 5, 0, 1, 1, 2, 3, 4}, paaBin)

	paa, _ = NewPaa(0, PdnTypeIPv6, nil, net.ParseIP("2001:db8::68"))
	paaBin = paa.Marshal()
	assert.Equal(t, []byte{0x4f, 0, 17, 0, 2,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x68,
	}, paaBin)

	paa, _ = NewPaa(0, PdnTypeIPv4v6, net.IPv4(1, 2, 3, 4), net.ParseIP("2001:db8::68"))
	paaBin = paa.Marshal()
	assert.Equal(t, []byte{0x4f, 0, 21, 0, 3,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x68,
		1, 2, 3, 4,
	}, paaBin)
}

func TestUnmarshal_Paa(t *testing.T) {
	paa, _ := NewPaa(0, PdnTypeIPv4, net.IPv4(1, 2, 3, 4), nil)
	paaBin := paa.Marshal()
	msg, tail, err := Unmarshal(paaBin)
	paa = msg.(*Paa)
	assert.Equal(t, byte(0), paa.instance)
	assert.Equal(t, PdnTypeIPv4, paa.Value)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), paa.ipv4)
	assert.Nil(t, paa.ipv6)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	paa, _ = NewPaa(1, PdnTypeIPv6, nil, net.ParseIP("2001:db8::68"))
	paaBin = paa.Marshal()
	msg, tail, err = Unmarshal(paaBin)
	paa = msg.(*Paa)
	assert.Equal(t, byte(1), paa.instance)
	assert.Equal(t, PdnTypeIPv6, paa.Value)
	assert.Nil(t, paa.ipv4)
	assert.Equal(t, net.ParseIP("2001:db8::68"), paa.ipv6)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	paa, _ = NewPaa(2, PdnTypeIPv4v6, net.IPv4(1, 2, 3, 4), net.ParseIP("2001:db8::68"))
	paaBin = paa.Marshal()
	msg, tail, err = Unmarshal(paaBin)
	paa = msg.(*Paa)
	assert.Equal(t, byte(2), paa.instance)
	assert.Equal(t, PdnTypeIPv4v6, paa.Value)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), paa.ipv4)
	assert.Equal(t, net.ParseIP("2001:db8::68"), paa.ipv6)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
