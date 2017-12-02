package domain

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSPgwRepo_AddGet(t *testing.T) {
	// init
	sgwCtrl := TheSgwCtrlRepo().GetCtrl(defaultSgwCtrlAddr)
	assert.Equal(t, defaultSgwCtrlAddr, sgwCtrl.UDPAddr())

	// Add & Get
	sgwCtrl2, err := newSgwCtrl(net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: GtpControlPort + 1,
	}, GtpUserPort+1, 1)
	assert.NoError(t, err)
	err = TheSgwCtrlRepo().AddCtrl(sgwCtrl2)
	assert.NoError(t, err)
	assert.Equal(t, sgwCtrl2, theSgwCtrlRepo.GetCtrl(sgwCtrl2.UDPAddr()))

	// Add error  :  duplication of SgwCtrlPort
	_, err = newSgwCtrl(net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: GtpControlPort,
	}, GtpUserPort+2, 1)
	assert.Error(t, err)

	// Add error : duplication of SgwDataPort
	_, err = newSgwCtrl(net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: GtpControlPort + 2,
	}, GtpUserPort, 1)
	assert.Error(t, err)
}
