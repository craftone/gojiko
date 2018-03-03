package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_LogLevel(t *testing.T) {
	// test sample config value
	assert.Equal(t, "debug", LogLevel("domain"))
}
