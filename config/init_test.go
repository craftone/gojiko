package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	Init()
	code := m.Run()
	os.Exit(code)
}

func TestInit_Gtpv2cTimeout(t *testing.T) {
	// test sample config value
	assert.Equal(t, 1000*time.Millisecond, Gtpv2cTimeoutDuration())

	// test update
	SetGtpv2cTimeout(3000)
	assert.Equal(t, 3000, Gtpv2cTimeout())
	assert.Equal(t, 3000*time.Millisecond, Gtpv2cTimeoutDuration())
}

func TestInit_Gtpv2cRetry(t *testing.T) {
	// test sample config value
	assert.Equal(t, 2, Gtpv2cRetry())

	// test update
	SetGtpv2cRetry(5)
	assert.Equal(t, 5, Gtpv2cRetry())
}
