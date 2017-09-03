package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecovery_new(t *testing.T) {
	rec := NewRecovery(1, 0)
	assert.Equal(t, recoveryNum, rec.header.typeNum)
	assert.Equal(t, byte(1), rec.Value)
}

func TestRecovery_marshal(t *testing.T) {
	var rec []byte
	rec = NewRecovery(0, 0).Marshal()
	assert.Equal(t, []byte{3, 0, 1, 0, 0}, rec)

	rec = NewRecovery(255, 0xf).Marshal()
	assert.Equal(t, []byte{3, 0, 1, 0, 255}, rec)
}
