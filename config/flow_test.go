package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_loadFlowConf(t *testing.T) {
	assert.Equal(t, 1*time.Millisecond, FlowUdpEchoWaitReceive())
}
