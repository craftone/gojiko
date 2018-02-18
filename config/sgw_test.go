package config

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_SGWs(t *testing.T) {
	sgw1 := Sgw{"sgw1", net.ParseIP("127.0.0.1")}
	assert.Equal(t, sgw1, GetSGWs()[0])
}
