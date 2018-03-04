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

	sendPackets, ok := fs.ReadInt64(SendPackets)
	assert.True(t, ok)
	assert.Equal(t, int64(1), sendPackets)
	sendBytes, ok := fs.ReadInt64(SendBytes)
	assert.True(t, ok)
	assert.Equal(t, int64(1460), sendBytes)
	_, ok = fs.ReadInt64(RecvPackets)
	assert.False(t, ok)

	// SetInt64 test
	fs.SetInt64(SendPackets, 0)
	sendPackets, ok = fs.ReadInt64(SendPackets)
	assert.True(t, ok)
	assert.Equal(t, int64(0), sendPackets)

	// cancel test
	cancel()
}
