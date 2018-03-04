package stats

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_FlowStats(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	fs := NewFlowStats(ctx)
	fs.SendInt64Msg(SendPackets, 1)
	fs.SendInt64Msg(SendBytes, 1460)

	// wait untill the receiver take all message
	for len(fs.toMsgReceiverChan) != 0 {
		time.Sleep(1 * time.Millisecond)
	}

	assert.Equal(t, int64(1), fs.ReadInt64(SendPackets))
	assert.Equal(t, int64(1460), fs.ReadInt64(SendBytes))
	assert.Equal(t, int64(0), fs.ReadInt64(RecvPackets))

	// SetInt64 test
	fs.SetInt64(SendPackets, 0)
	assert.Equal(t, int64(0), fs.ReadInt64(SendPackets))

	// cancel test
	cancel()
}
