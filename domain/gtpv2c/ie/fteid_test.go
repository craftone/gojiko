package ie

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/stretchr/testify/assert"
)

func TestNewFteid(t *testing.T) {
	fteid, err := NewFteid(1, net.IPv4(1, 2, 3, 4), nil, S5S8SgwGtpCIf, 0x0006c6ea)
	assert.Equal(t, fteidNum, fteid.typeNum)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), fteid.Ipv4())
	assert.Nil(t, fteid.Ipv6())
	assert.Equal(t, S5S8SgwGtpCIf, fteid.IfType())
	assert.Equal(t, gtp.Teid(0x0006c6ea), fteid.Teid())
	assert.Nil(t, err)
}

func TestFteid_Marshal(t *testing.T) {
	fteid, _ := NewFteid(1, net.IPv4(1, 2, 3, 4), nil, 6, 0x0006c6ea)
	fteidBin := fteid.Marshal()
	assert.Equal(t, []byte{0x57, 0, 9, 1, 0x86, 0, 0x06, 0xc6, 0xea, 1, 2, 3, 4}, fteidBin)
}

func TestUnmarshal_fteid(t *testing.T) {
	fteid, _ := NewFteid(1, net.IPv4(1, 2, 3, 4), nil, S5S8SgwGtpCIf, 0x0006c6ea)
	fteidBin := fteid.Marshal()
	msg, tail, err := Unmarshal(fteidBin, CreateSessionRequest)
	fteid = msg.(*Fteid)
	assert.Equal(t, fteidNum, fteid.typeNum)
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), fteid.Ipv4())
	assert.Nil(t, fteid.Ipv6())
	assert.Equal(t, S5S8SgwGtpCIf, fteid.IfType())
	assert.Equal(t, gtp.Teid(0x0006c6ea), fteid.Teid())
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
