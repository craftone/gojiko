package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_unknown(t *testing.T) {
	rec, err := NewRecovery(0, 255)
	recBin := rec.Marshal()
	recBin[0] = 0
	msg, _, err := Unmarshal(recBin, MsToNetwork)
	assert.Nil(t, msg)
	assert.Error(t, err)
}
