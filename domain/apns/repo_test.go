package apns

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApnRepo_PostFind(t *testing.T) {
	// Create first APN
	ips := []net.IP{net.IPv4(127, 0, 0, 1)}
	apn, err := NewApn("apn.example.com", "440", "10", ips)
	assert.NoError(t, err)

	err = TheRepo().Post(apn)
	assert.NoError(t, err)

	// Create second APN
	apn2, _ := NewApn("apn.example2.com", "440", "10", ips)
	err = TheRepo().Post(apn2)
	assert.NoError(t, err)

	// Create dupulication error
	apnDup, _ := NewApn("APN.Example.Com", "440", "10", ips)
	err = TheRepo().Post(apnDup)
	assert.EqualError(t, err, "There is already the name's APN : apn.example.com.mnc010.mcc440.gprs")

	// Find success
	apnFinded, err := TheRepo().Find("apn.example.com", "440", "10")
	assert.Equal(t, apn, apnFinded)

	// Find unsuccess
	apnFinded, err = TheRepo().Find("noexist.example.com", "440", "10")
	assert.Nil(t, apnFinded)
	assert.EqualError(t, err, "There is no such APN : networkID=noexist.example.com MCC=440 MNC=10")
}
