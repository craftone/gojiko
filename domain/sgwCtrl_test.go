package domain

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/domain/apns"
	"github.com/stretchr/testify/assert"
)

func TestSgwCtrl_CreateSession(t *testing.T) {
	Init()

	apn, _ := apns.NewApn("example.com", "440", "10", []net.IP{net.IPv4(127, 1, 1, 1)})
	apns.TheRepo().Post(apn)

	sgwCtrl := theSgwCtrlRepo.getCtrl(defaultSgwCtrlAddr)
	err := sgwCtrl.CreateSession(
		"440101234567890", "819012345678", "0123456789012345",
		"440", "10", "example.com", 5,
	)
	assert.NoError(t, err)
}
