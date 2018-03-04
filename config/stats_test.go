package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadStatsConf(t *testing.T) {
	assert.Equal(t, 100, StatsFlowBufferLength())
}
