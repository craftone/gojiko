package apns

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApn(t *testing.T) {
	// normal
	ips := []net.IP{net.IPv4(127, 0, 0, 1)}
	apn, err := NewApn("apn.example.com", "440", "10", ips)
	assert.NoError(t, err)
	assert.Equal(t, apn.FullString(), "apn.example.com.mnc010.mcc440.gprs")

	// 3-digits mnc
	apn, err = NewApn("apn.example.com", "440", "110", ips)
	assert.NoError(t, err)
	assert.Equal(t, apn.FullString(), "apn.example.com.mnc110.mcc440.gprs")

	// Upper case network ID
	apn, err = NewApn("APN.Example.Com", "440", "110", ips)
	assert.NoError(t, err)
	assert.Equal(t, apn.FullString(), "apn.example.com.mnc110.mcc440.gprs")

	// invalid network id
	invalidNetworkIDs := []string{
		"apn.examp%le.com",
		"apn.too" + "oooooooooo" + "oooooooooo" + "oooooooooo" + "oooooooooo" + "oooooooooo" + "long.example.com",
		"apn.*.com",
		"rac.example.com",
		"lac.example.com",
		"sgsn.example.com",
		"rnc.example.com",
		"apn.example.com.gprs",
	}
	for _, nid := range invalidNetworkIDs {
		apn, err = NewApn(nid, "440", "10", ips)
		assert.Error(t, err)
	}

	// invalid mnc
	apn, err = NewApn("apn.example.com", "4400", "10", ips)
	assert.Error(t, err)
	apn, err = NewApn("apn.example.com", "40", "10", ips)
	assert.Error(t, err)
	apn, err = NewApn("apn.example.com", "44A", "10", ips)
	assert.Error(t, err)

	// invalid mcc
	apn, err = NewApn("apn.example.com", "440", "1000", ips)
	assert.Error(t, err)
	apn, err = NewApn("apn.example.com", "440", "1", ips)
	assert.Error(t, err)
	apn, err = NewApn("apn.example.com", "440", "10A", ips)
	assert.Error(t, err)

	// no ips
	apn, err = NewApn("apn.example.com", "440", "10", []net.IP{})
	assert.Error(t, err)

	apn, err = NewApn("apn.example.com", "440", "10", nil)
	assert.Error(t, err)
}

func TestGetIP(t *testing.T) {
	// 2 ips
	ips := []net.IP{net.IPv4(127, 0, 0, 1), net.IPv4(127, 0, 0, 2)}
	apn, _ := NewApn("apn.example.com", "440", "10", ips)

	ip := apn.GetIP()
	assert.Equal(t, ips[0], ip)
	ip = apn.GetIP()
	assert.Equal(t, ips[1], ip)
	ip = apn.GetIP()
	assert.Equal(t, ips[0], ip)
	ip = apn.GetIP()
	assert.Equal(t, ips[1], ip)

	// only 1 ip
	ips = []net.IP{net.IPv4(127, 0, 0, 1)}
	apn, _ = NewApn("apn.example.com", "440", "10", ips)
	ip = apn.GetIP()
	assert.Equal(t, ips[0], ip)
	ip = apn.GetIP()
	assert.Equal(t, ips[0], ip)
}
