package ie

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/gtpv2c/ie/pco"
	"github.com/stretchr/testify/assert"
)

func TestNewPcoMsToNetwork_marshal(t *testing.T) {
	ipcp := pco.NewIpcp(pco.ConfigureRequest, 0, net.IPv4(0, 0, 0, 0), net.IPv4(0, 0, 0, 0))
	msToNet := pco.NewMsToNetwork(ipcp, true, false, true)
	pcoMsToNetwork, _ := NewPcoMsToNetwork(0, msToNet)
	pcoMsToNetworkBin := pcoMsToNetwork.Marshal()
	assert.Equal(t, []byte{
		0x4E, 00, 0x1A, 00,
		0x80,       // PCO octet 3
		0x80, 0x21, // IPCP 8021H
		0x10,     // Length : 16
		0x01,     // Code : Configure-Request
		0x00,     // Identifier : 0
		0x00, 12, // Length : 12
		0x81, 6, 00, 00, 00, 00, // Option 129 Primary DNS
		0x83, 6, 00, 00, 00, 00, // Option 131 Secondary DNS
		00, 0x0D, 00, // DNS Server IPv4 Address Request 000DH
		00, 0x0A, 00, // IP address allocation via NAS Signalling 000AH
	}, pcoMsToNetworkBin)

}

func TestUnmarshal_PcoMsToNetwork(t *testing.T) {
	ipcp := pco.NewIpcp(pco.ConfigureRequest, 0, net.IPv4(0, 0, 0, 0), net.IPv4(0, 0, 0, 0))
	msToNet := pco.NewMsToNetwork(ipcp, true, false, true)
	pcoMsToNetwork, _ := NewPcoMsToNetwork(0, msToNet)
	pcoMsToNetworkBin := pcoMsToNetwork.Marshal()
	msg, tail, err := Unmarshal(pcoMsToNetworkBin, CreateSessionRequest)
	pcoMsToNetwork = msg.(*PcoMsToNetwork)
	assert.Equal(t, byte(0), pcoMsToNetwork.instance)
	assert.Equal(t, pco.ConfigureRequest, pcoMsToNetwork.Ipcp().Code())
	assert.Equal(t, byte(0), pcoMsToNetwork.Ipcp().Identifier())
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), pcoMsToNetwork.Ipcp().PriDNS())
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), pcoMsToNetwork.Ipcp().SecDNS())
	assert.Equal(t, true, pcoMsToNetwork.DNSServerV4Req())
	assert.Equal(t, false, pcoMsToNetwork.DNSServerV6Req())
	assert.Equal(t, true, pcoMsToNetwork.IPAllocViaNasSig())
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	//ensure no refference to the buffer
	copy(pcoMsToNetworkBin, make([]byte, len(pcoMsToNetworkBin)))
	assert.Equal(t, pco.ConfigureRequest, pcoMsToNetwork.Ipcp().Code())
	assert.Equal(t, byte(0), pcoMsToNetwork.Ipcp().Identifier())
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), pcoMsToNetwork.Ipcp().PriDNS())
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), pcoMsToNetwork.Ipcp().SecDNS())
	assert.Equal(t, true, pcoMsToNetwork.DNSServerV4Req())
	assert.Equal(t, false, pcoMsToNetwork.DNSServerV6Req())
	assert.Equal(t, true, pcoMsToNetwork.IPAllocViaNasSig())
}
