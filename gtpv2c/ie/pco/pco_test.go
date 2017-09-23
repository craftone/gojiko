package pco

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPco_Marshal(t *testing.T) {
	// DnsServerV4 only
	p := Pco{
		ConfigProto:  0,
		DnsServerV4s: []*DnsServerV4{NewDnsServerV4(net.IPv4(1, 2, 3, 4))},
	}
	pBin := p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0d, 4, 1, 2, 3, 4}, pBin)

	// DnsServerV4 *2 only
	p = Pco{
		ConfigProto: 0,
		DnsServerV4s: []*DnsServerV4{
			NewDnsServerV4(net.IPv4(1, 2, 3, 4)),
			NewDnsServerV4(net.IPv4(5, 6, 7, 8)),
		},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80,
		0x00, 0x0d, 4, 1, 2, 3, 4,
		0x00, 0x0d, 4, 5, 6, 7, 8,
	}, pBin)

	// DnsServerV6 only
	p = Pco{
		ConfigProto:  0,
		DnsServerV6s: []*DnsServerV6{NewDnsServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x03, 16,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0x68,
	}, pBin)

	// DnsServerV4 and DnsServerV6
	p = Pco{
		ConfigProto:  0,
		DnsServerV4s: []*DnsServerV4{NewDnsServerV4(net.IPv4(1, 2, 3, 4))},
		DnsServerV6s: []*DnsServerV6{NewDnsServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80,
		0x00, 0x0d, 4, 1, 2, 3, 4,
		0x00, 0x03, 16,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0x68,
	}, pBin)
}

func TestUnmarshal(t *testing.T) {
	// DnsServerV4 only
	p := Pco{
		ConfigProto:  0,
		DnsServerV4s: []*DnsServerV4{NewDnsServerV4(net.IPv4(1, 2, 3, 4))},
	}
	pBin := p.Marshal()
	p, tail, err := Unmarshal(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DnsServerV4s[0].value)
	assert.Equal(t, 0, len(p.DnsServerV6s))
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV6 only
	p = Pco{
		ConfigProto:  0,
		DnsServerV6s: []*DnsServerV6{NewDnsServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	p, tail, err = Unmarshal(pBin)
	assert.Equal(t, 0, len(p.DnsServerV4s))
	assert.Equal(t, net.ParseIP("2001:db8::68"), p.DnsServerV6s[0].value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV4 *2 and DnsServerV6 * 2
	p = Pco{
		ConfigProto: 0,
		DnsServerV4s: []*DnsServerV4{
			NewDnsServerV4(net.IPv4(1, 2, 3, 4)),
			NewDnsServerV4(net.IPv4(5, 6, 7, 8)),
		},
		DnsServerV6s: []*DnsServerV6{
			NewDnsServerV6(net.ParseIP("2001:db8::68")),
			NewDnsServerV6(net.ParseIP("2001:db8::69")),
		},
	}
	pBin = p.Marshal()
	p, tail, err = Unmarshal(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DnsServerV4s[0].value)
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), p.DnsServerV4s[1].value)
	assert.Equal(t, net.ParseIP("2001:db8::68"), p.DnsServerV6s[0].value)
	assert.Equal(t, net.ParseIP("2001:db8::69"), p.DnsServerV6s[1].value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
