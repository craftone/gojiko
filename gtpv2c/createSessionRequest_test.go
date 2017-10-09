package gtpv2c

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCreateSessionRequest(t *testing.T) {
	// imsi, _ := ie.NewImsi(0, "23434")
	csReqArg := CreateSessionRequestArg{
	// Imsi: imsi,
	}
	_, err := NewCreateSessionRequest(0, csReqArg)
	assert.Error(t, err)
}
