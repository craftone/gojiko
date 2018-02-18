package config

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_APNs(t *testing.T) {
	apn1 := Apn{"apn-example.com", "440", "10", []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.2")}}
	assert.Equal(t, apn1, GetAPNs()[0])
	apn2 := Apn{"apn-example2.com", "440", "10", []net.IP{net.ParseIP("127.0.2.1"), net.ParseIP("127.0.2.2")}}
	assert.Equal(t, apn2, GetAPNs()[1])
}
