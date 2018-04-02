package config

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_SGWs(t *testing.T) {
	sgw1 := Sgw{"sgw1", net.ParseIP("127.0.0.1"), byte(1)}
	fmt.Printf("%#v\n", GetSGWs()[0])
	assert.Equal(t, sgw1, GetSGWs()[0])
}
