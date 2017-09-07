package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_unknown(t *testing.T) {
	recBin := NewRecovery(0, 255).Marshal()
	recBin[0] = 0
	msg, _, err := Unmarshal(recBin)
	assert.Nil(t, msg)
	assert.Error(t, err)
}
