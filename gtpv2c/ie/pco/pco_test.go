package pco

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPco_Marshal(t *testing.T) {
	// DnsServerV4 only
	p := Pco{
		ConfigProto: 0,
		DnsServerV4: NewDnsServerV4(net.IPv4(1, 2, 3, 4)),
	}
	pBin := p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0d, 4, 1, 2, 3, 4}, pBin)

	// DnsServerV6 only
	p = Pco{
		ConfigProto: 0,
		DnsServerV6: NewDnsServerV6(net.ParseIP("2001:db8::68")),
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x03, 16,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0x68,
	}, pBin)

	// DnsServerV4 and DnsServerV6
	p = Pco{
		ConfigProto: 0,
		DnsServerV4: NewDnsServerV4(net.IPv4(1, 2, 3, 4)),
		DnsServerV6: NewDnsServerV6(net.ParseIP("2001:db8::68")),
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0,
		0x00, 0x0d, 4, 1, 2, 3, 4,
		0x00, 0x03, 16,
		0, 0, 0, 0, 0, 0, 0, 0x68,
	}, pBin)
}

func TestUnmarshal(t *testing.T) {
	// DnsServerV4 only
	p := Pco{
		ConfigProto: 0,
		DnsServerV4: NewDnsServerV4(net.IPv4(1, 2, 3, 4)),
	}
	pBin := p.Marshal()
	p, tail, err := Unmarshal(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DnsServerV4.value)
	assert.Nil(t, p.DnsServerV6)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV6 only
	p = Pco{
		ConfigProto: 0,
		DnsServerV6: NewDnsServerV6(net.ParseIP("2001:db8::68")),
	}
	pBin = p.Marshal()
	p, tail, err = Unmarshal(pBin)
	assert.Nil(t, p.DnsServerV4)
	assert.Equal(t, net.ParseIP("2001:db8::68"), p.DnsServerV6.value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV4 and DnsServerV6
	p = Pco{
		ConfigProto: 0,
		DnsServerV4: NewDnsServerV4(net.IPv4(1, 2, 3, 4)),
		DnsServerV6: NewDnsServerV6(net.ParseIP("2001:db8::68")),
	}
	pBin = p.Marshal()
	p, tail, err = Unmarshal(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DnsServerV4.value)
	assert.Equal(t, net.ParseIP("2001:db8::68"), p.DnsServerV6.value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
