package pco

import (
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPcoMsToNetwork_Marshal(t *testing.T) {
	// DnsServerV4 only
	p := PcoMsToNetwork{
		pco:            pco{},
		DnsServerV4Req: true,
	}
	pBin := p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0d, 0}, pBin)

	// DnsServerV6 only
	p = PcoMsToNetwork{
		pco:            pco{},
		DnsServerV6Req: true,
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x03, 0}, pBin)

	// DnsServerV4 and DnsServerV6
	p = PcoMsToNetwork{
		pco:            pco{},
		DnsServerV4Req: true,
		DnsServerV6Req: true,
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0d, 0, 0x00, 0x03, 0}, pBin)
}

func TestPcoNetworkToMs_Marshal(t *testing.T) {
	// DnsServerV4 only
	p := PcoNetworkToMs{
		pco:          pco{},
		DnsServerV4s: []*DnsServerV4{NewDnsServerV4(net.IPv4(1, 2, 3, 4))},
	}
	pBin := p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x0d, 4, 1, 2, 3, 4}, pBin)

	// DnsServerV4 *2 only
	p = PcoNetworkToMs{
		pco: pco{},
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
	p = PcoNetworkToMs{
		pco:          pco{},
		DnsServerV6s: []*DnsServerV6{NewDnsServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	assert.Equal(t, []byte{0x80, 0x00, 0x03, 16,
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0x68,
	}, pBin)

	// DnsServerV4 and DnsServerV6
	p = PcoNetworkToMs{
		pco:          pco{},
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

func Test_pco_marshal(t *testing.T) {
	type fields struct {
		ConfigProto byte
	}
	type args struct {
		body []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pco{
				ConfigProto: tt.fields.ConfigProto,
			}
			if got := p.marshal(tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pco.marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalMsToNetowrk(t *testing.T) {
	// DnsServerV4 only
	p := PcoMsToNetwork{
		pco:            pco{},
		DnsServerV4Req: true,
	}
	pBin := p.Marshal()
	p, tail, err := UnmarshalMsToNetowrk(pBin)
	assert.Equal(t, true, p.DnsServerV4Req)
	assert.Equal(t, false, p.DnsServerV6Req)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV6 only
	p = PcoMsToNetwork{
		pco:            pco{},
		DnsServerV6Req: true,
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalMsToNetowrk(pBin)
	assert.Equal(t, false, p.DnsServerV4Req)
	assert.Equal(t, true, p.DnsServerV6Req)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV4 and DnsServerV6
	p = PcoMsToNetwork{
		pco:            pco{},
		DnsServerV4Req: true,
		DnsServerV6Req: true,
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalMsToNetowrk(pBin)
	assert.Equal(t, true, p.DnsServerV4Req)
	assert.Equal(t, true, p.DnsServerV6Req)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshalNetowrkToMs(t *testing.T) {
	// DnsServerV4 only
	p := PcoNetworkToMs{
		pco:          pco{},
		DnsServerV4s: []*DnsServerV4{NewDnsServerV4(net.IPv4(1, 2, 3, 4))},
	}
	pBin := p.Marshal()
	p, tail, err := UnmarshalNetowrkToMs(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DnsServerV4s[0].value)
	assert.Nil(t, p.DnsServerV6s)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV4 *2 only
	p = PcoNetworkToMs{
		pco: pco{},
		DnsServerV4s: []*DnsServerV4{
			NewDnsServerV4(net.IPv4(1, 2, 3, 4)),
			NewDnsServerV4(net.IPv4(5, 6, 7, 8)),
		},
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalNetowrkToMs(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DnsServerV4s[0].value)
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), p.DnsServerV4s[1].value)
	assert.Nil(t, p.DnsServerV6s)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV6 only
	p = PcoNetworkToMs{
		pco:          pco{},
		DnsServerV6s: []*DnsServerV6{NewDnsServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalNetowrkToMs(pBin)
	assert.Nil(t, p.DnsServerV4s)
	assert.Equal(t, net.ParseIP("2001:db8::68"), p.DnsServerV6s[0].value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	// DnsServerV4 and DnsServerV6
	p = PcoNetworkToMs{
		pco:          pco{},
		DnsServerV4s: []*DnsServerV4{NewDnsServerV4(net.IPv4(1, 2, 3, 4))},
		DnsServerV6s: []*DnsServerV6{NewDnsServerV6(net.ParseIP("2001:db8::68"))},
	}
	pBin = p.Marshal()
	p, tail, err = UnmarshalNetowrkToMs(pBin)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), p.DnsServerV4s[0].value)
	assert.Equal(t, net.ParseIP("2001:db8::68"), p.DnsServerV6s[0].value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
