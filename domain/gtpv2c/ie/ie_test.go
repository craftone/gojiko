package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_unknown(t *testing.T) {
	rec, err := NewRecovery(0, 255)
	recBin := rec.Marshal()
	recBin[0] = 255
	msg, _, err := Unmarshal(recBin, CreateSessionRequest)
	assert.Nil(t, msg)
	assert.Equal(t, err, &UnknownIEError{255, 0})
}
