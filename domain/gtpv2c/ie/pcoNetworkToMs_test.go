package ie

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/domain/gtpv2c/ie/pco"
	"github.com/stretchr/testify/assert"
)

func TestNewPcoNetworkToMs_marshal(t *testing.T) {
	ipcp := pco.NewIpcp(pco.ConfigureNack, 0, net.IPv4(1, 2, 3, 4), net.IPv4(5, 6, 7, 8))
	dnsServerV4s := []*pco.DNSServerV4{
		pco.NewDNSServerV4(net.IPv4(1, 2, 3, 4)),
		pco.NewDNSServerV4(net.IPv4(5, 6, 7, 8)),
	}
	dnsServerV6s := []*pco.DNSServerV6{}
	netToMs := pco.NewNetworkToMs(ipcp, dnsServerV4s, dnsServerV6s)
	pcoNetworkToMs, _ := NewPcoNetworkToMs(0, netToMs)
	pcoNetworkToMsBin := pcoNetworkToMs.Marshal()
	assert.Equal(t, []byte{
		0x4E, 00, 0x22, 00,
		0x80,       // PCO octet 3
		0x80, 0x21, // IPCP 8021H
		0x10,     // Length : 16
		0x03,     // Code : Configure-Nack
		0x00,     // Identifier : 0
		0x00, 12, // Length : 12
		0x81, 6, 1, 2, 3, 4, // Option 129 Primary DNS
		0x83, 6, 5, 6, 7, 8, // Option 131 Secondary DNS
		00, 0x0D, 4, 1, 2, 3, 4, // DNS Server IPv4 Address Request 000DH
		00, 0x0D, 4, 5, 6, 7, 8, // DNS Server IPv4 Address Request 000DH
	}, pcoNetworkToMsBin)

}

func TestUnmarshal_PcoNetworkToMs(t *testing.T) {
	ipcp := pco.NewIpcp(pco.ConfigureNack, 0, net.IPv4(1, 2, 3, 4), net.IPv4(5, 6, 7, 8))
	dnsServerV4s := []*pco.DNSServerV4{
		pco.NewDNSServerV4(net.IPv4(1, 2, 3, 4)),
		pco.NewDNSServerV4(net.IPv4(5, 6, 7, 8)),
	}
	dnsServerV6s := []*pco.DNSServerV6{}
	netToMs := pco.NewNetworkToMs(ipcp, dnsServerV4s, dnsServerV6s)
	pcoNetworkToMs, _ := NewPcoNetworkToMs(0, netToMs)
	pcoNetworkToMsBin := pcoNetworkToMs.Marshal()
	msg, tail, err := Unmarshal(pcoNetworkToMsBin, CreateSessionResponse)
	pcoNetworkToMs = msg.(*PcoNetworkToMs)
	assert.Equal(t, byte(0), pcoNetworkToMs.instance)
	assert.Equal(t, pco.ConfigureNack, pcoNetworkToMs.Ipcp().Code())
	assert.Equal(t, byte(0), pcoNetworkToMs.Ipcp().Identifier())
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), pcoNetworkToMs.Ipcp().PriDNS())
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), pcoNetworkToMs.Ipcp().SecDNS())
	assert.Equal(t, pco.NewDNSServerV4(net.IPv4(1, 2, 3, 4)), pcoNetworkToMs.DNSServerV4s()[0])
	assert.Equal(t, pco.NewDNSServerV4(net.IPv4(5, 6, 7, 8)), pcoNetworkToMs.DNSServerV4s()[1])
	assert.Equal(t, 0, len(pcoNetworkToMs.DNSServerV6s()))
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	//ensure no refference to the buffer
	copy(pcoNetworkToMsBin, make([]byte, len(pcoNetworkToMsBin)))
	assert.Equal(t, pco.ConfigureNack, pcoNetworkToMs.Ipcp().Code())
	assert.Equal(t, byte(0), pcoNetworkToMs.Ipcp().Identifier())
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), pcoNetworkToMs.Ipcp().PriDNS())
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), pcoNetworkToMs.Ipcp().SecDNS())
	assert.Equal(t, pco.NewDNSServerV4(net.IPv4(1, 2, 3, 4)), pcoNetworkToMs.DNSServerV4s()[0])
	assert.Equal(t, pco.NewDNSServerV4(net.IPv4(5, 6, 7, 8)), pcoNetworkToMs.DNSServerV4s()[1])
	assert.Equal(t, 0, len(pcoNetworkToMs.DNSServerV6s()))
}
