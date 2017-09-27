package pco

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPcoMsToNetwork_Marshal(t *testing.T) {
	// DNSServerV4 only
	p := MsToNetwork{
		pco:            pco{},
		DNSServerV4Req: true,
	}
	pBin := p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0d, 0}, pBin)

	// DNSServerV6 only
	p = MsToNetwork{
		pco:            pco{},
		DNSServerV6Req: true,
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x03, 0}, pBin)

	// DNSServerV4 and DNSServerV6
	p = MsToNetwork{
		pco:            pco{},
		DNSServerV4Req: true,
		DNSServerV6Req: true,
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0d, 0, 0x00, 0x03, 0}, pBin)

	// IPCP only
	p = MsToNetwork{
		pco: pco{
			Ipcp: NewIpcp(ConfigureRequest, 0, net.IPv4(0, 0, 0, 0), net.IPv4(0, 0, 0, 0)),
		},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{
		0x80,       // PCO octet 3
		0x80, 0x21, // IPCP 8021H
		0x10,     // Length : 16
		1,        // Code : Configure-Request
		0,        // Identifier : 0
		00, 0x0c, // Length: 12
		0x81,       // Option : 129 Primary DNS
		6,          // Length : 6
		0, 0, 0, 0, // 0.0.0.0
		0x83,       // Option : 131 Secondary DNS
		6,          // Length : 6
		0, 0, 0, 0, // 0.0.0.0
	}, pBin)

	// IPAddrAllocSig only
	p = MsToNetwork{
		pco:              pco{},
		IPAllocViaNasSig: true,
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0a, 0}, pBin)
}

func TestPcoNetworkToMs_Marshal(t *testing.T) {
	// DNSServerV4 only
	p := NetworkToMs{
		pco:          pco{},
		DNSServerV4s: []*DNSServerV4{NewDNSServerV4(net.IPv4(1, 2, 3, 4))},
	}
	pBin := p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0d, 4, 1, 2, 3, 4}, pBin)

	// DNSServerV4 *2 only
	p = NetworkToMs{
		pco: pco{},
		DNSServerV4s: []*DNSServerV4{
			NewDNSServerV4(net.IPv4(1, 2, 3, 4)),
			NewDNSServerV4(net.IPv4(5, 6, 7, 8)),
		},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80,
		0x00, 0x0d, 4, 1, 2, 3, 4,
		0x00, 0x0d, 4, 5, 6, 7, 8,
	}, pBin)

	// DNSServerV6 only
	p = NetworkToMs{
		pco:          pco{},
		DNSServerV6s: []*DNSServerV6{NewDNSServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x03, 16,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0x68,
	}, pBin)

	// DNSServerV4 and DNSServerV6
	p = NetworkToMs{
		pco:          pco{},
		DNSServerV4s: []*DNSServerV4{NewDNSServerV4(net.IPv4(1, 2, 3, 4))},
		DNSServerV6s: []*DNSServerV6{NewDNSServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80,
		0x00, 0x0d, 4, 1, 2, 3, 4,
		0x00, 0x03, 16,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0x68,
	}, pBin)

	// IPCP only
	p = NetworkToMs{
		pco: pco{
			Ipcp: NewIpcp(ConfigureRequest, 0, net.IPv4(1, 2, 3, 4), net.IPv4(5, 6, 7, 8)),
		},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{
		0x80,       // PCO octet 3
		0x80, 0x21, // IPCP 8021H
		0x10,     // Length : 16
		1,        // Code : Configure-Request
		0,        // Identifier : 0
		00, 0x0c, // Length: 12
		0x81,       // Option : 129 Primary DNS
		6,          // Length : 6
		1, 2, 3, 4, // 1.2.3.4
		0x83,       // Option : 131 Secondary DNS
		6,          // Length : 6
		5, 6, 7, 8, // 5.6.7.8
	}, pBin)

}

func TestUnmarshalMsToNetowrk(t *testing.T) {
	// DNSServerV4 only
	p := &MsToNetwork{
		pco:            pco{},
		DNSServerV4Req: true,
	}
	pBin := p.Marshal()
	p, tail, err := UnmarshalMsToNetowrk(pBin)
	assert.Equal(t, true, p.DNSServerV4Req)
	assert.Equal(t, false, p.DNSServerV6Req)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DNSServerV6 only
	p = &MsToNetwork{
		pco:            pco{},
		DNSServerV6Req: true,
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalMsToNetowrk(pBin)
	assert.Equal(t, false, p.DNSServerV4Req)
	assert.Equal(t, true, p.DNSServerV6Req)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DNSServerV4 and DNSServerV6
	p = &MsToNetwork{
		pco:            pco{},
		DNSServerV4Req: true,
		DNSServerV6Req: true,
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalMsToNetowrk(pBin)
	assert.Equal(t, true, p.DNSServerV4Req)
	assert.Equal(t, true, p.DNSServerV6Req)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// IPCP only
	p = &MsToNetwork{
		pco: pco{
			Ipcp: NewIpcp(ConfigureRequest, 0, net.IPv4(0, 0, 0, 0), net.IPv4(0, 0, 0, 0)),
		},
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalMsToNetowrk(pBin)
	assert.Equal(t, false, p.DNSServerV4Req)
	assert.Equal(t, false, p.DNSServerV6Req)
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), p.Ipcp.PriDNS)
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), p.Ipcp.SecDNS)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// IPAddrAllocSig only
	p = &MsToNetwork{
		pco:              pco{},
		IPAllocViaNasSig: true,
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalMsToNetowrk(pBin)
	assert.Equal(t, false, p.DNSServerV4Req)
	assert.Equal(t, false, p.DNSServerV6Req)
	assert.Equal(t, true, p.IPAllocViaNasSig)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshalNetowrkToMs(t *testing.T) {
	// DNSServerV4 only
	p := &NetworkToMs{
		pco:          pco{},
		DNSServerV4s: []*DNSServerV4{NewDNSServerV4(net.IPv4(1, 2, 3, 4))},
	}
	pBin := p.Marshal()
	p, tail, err := UnmarshalNetowrkToMs(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DNSServerV4s[0].value)
	assert.Nil(t, p.DNSServerV6s)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DNSServerV4 *2 only
	p = &NetworkToMs{
		pco: pco{},
		DNSServerV4s: []*DNSServerV4{
			NewDNSServerV4(net.IPv4(1, 2, 3, 4)),
			NewDNSServerV4(net.IPv4(5, 6, 7, 8)),
		},
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalNetowrkToMs(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DNSServerV4s[0].value)
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), p.DNSServerV4s[1].value)
	assert.Nil(t, p.DNSServerV6s)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DNSServerV6 only
	p = &NetworkToMs{
		pco:          pco{},
		DNSServerV6s: []*DNSServerV6{NewDNSServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalNetowrkToMs(pBin)
	assert.Nil(t, p.DNSServerV4s)
	assert.Equal(t, net.ParseIP("2001:db8::68"), p.DNSServerV6s[0].value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DNSServerV4 and DNSServerV6
	p = &NetworkToMs{
		pco:          pco{},
		DNSServerV4s: []*DNSServerV4{NewDNSServerV4(net.IPv4(1, 2, 3, 4))},
		DNSServerV6s: []*DNSServerV6{NewDNSServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalNetowrkToMs(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DNSServerV4s[0].value)
	assert.Equal(t, net.ParseIP("2001:db8::68"), p.DNSServerV6s[0].value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// IPCP only
	p = &NetworkToMs{
		pco: pco{
			Ipcp: NewIpcp(ConfigureRequest, 0, net.IPv4(0, 0, 0, 0), net.IPv4(0, 0, 0, 0)),
		},
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalNetowrkToMs(pBin)
	assert.Nil(t, p.DNSServerV4s)
	assert.Nil(t, p.DNSServerV6s)
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), p.Ipcp.PriDNS)
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), p.Ipcp.SecDNS)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
