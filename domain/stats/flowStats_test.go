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
	fs.SendUint64Msg(SendPackets, 1)
	fs.SendUint64Msg(SendBytes, 1460)

	// wait untill the receiver take all message
	for len(fs.toMsgReceiverChan) != 0 {
		time.Sleep(1 * time.Millisecond)
	}

	assert.Equal(t, uint64(1), fs.ReadUint64(SendPackets))
	assert.Equal(t, uint64(1460), fs.ReadUint64(SendBytes))
	assert.Equal(t, uint64(0), fs.ReadUint64(RecvPackets))

	// SetInt64 test
	fs.SetUint64(SendPackets, 0)
	assert.Equal(t, uint64(0), fs.ReadUint64(SendPackets))

	// cancel test
	cancel()
}

func Test_TimeMsg(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	fs := NewFlowStats(ctx)
	now1 := time.Now()
	fs.SendTimeMsg(StartTime, now1)
	now2 := time.Now()
	fs.SendTimeMsg(EndTime, now2)

	// wait untill the receiver take all message
	for len(fs.toMsgReceiverChan) != 0 {
		time.Sleep(1 * time.Millisecond)
	}

	assert.Equal(t, now1, fs.ReadTime(StartTime))
	assert.Equal(t, now2, fs.ReadTime(EndTime))
}

func Test_AllItems(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	fs := NewFlowStats(ctx)
	startTime := time.Now()
	fs.SendTimeMsg(StartTime, startTime)
	fs.SendUint64Msg(SendBytes, 100000000)
	fs.SendUint64Msg(SendPackets, 100)
	fs.SendUint64Msg(RecvBytes, 1000000000)
	fs.SendUint64Msg(RecvPackets, 10000)
	fs.SendUint64Msg(SendBytesSkipped, 1024)
	fs.SendUint64Msg(SendPacketsSkipped, 1000)
	fs.SendUint64Msg(RecvBytesInvalid, 10240)
	fs.SendUint64Msg(RecvPacketsInvalid, 10000)
	endTime := startTime.Add(3 * time.Second)
	fs.SendTimeMsg(EndTime, endTime)

	// wait untill the receiver take all message
	for len(fs.toMsgReceiverChan) != 0 {
		time.Sleep(1 * time.Millisecond)
	}

	// assert bitrates
	sendBitrate, sendBitrateStr := fs.SendBitrate()
	assert.Equal(t, float64(100000000)*8/3, sendBitrate)
	assert.Equal(t, "266.7 Mbps", sendBitrateStr)
	recvBitrate, recvBitrateStr := fs.RecvBitrate()
	assert.Equal(t, float64(1000000000)*8/3, recvBitrate)
	assert.Equal(t, "2.7 Gbps", recvBitrateStr)

	// assert normal bytes, packets
	_, sendBytes := fs.SendBytes()
	assert.Equal(t, "95 MiB", sendBytes)
	_, recvBytes := fs.RecvBytes()
	assert.Equal(t, "954 MiB", recvBytes)
	_, sendPkts := fs.SendPackets()
	assert.Equal(t, "100 pkts", sendPkts)
	_, recvPkts := fs.RecvPackets()
	assert.Equal(t, "10.0 kpkts", recvPkts)

	// assert skipped & invalid bytes, packets
	_, sendBytesSkipped := fs.SendBytesSkipped()
	assert.Equal(t, "1.0 KiB", sendBytesSkipped)
	_, recvBytesInvalid := fs.RecvBytesInvalid()
	assert.Equal(t, "10 KiB", recvBytesInvalid)
	_, sendPktsSkipped := fs.SendPacketsSkipped()
	assert.Equal(t, "1.0 kpkts", sendPktsSkipped)
	_, recvPktsInvalid := fs.RecvPacketsInvalid()
	assert.Equal(t, "10.0 kpkts", recvPktsInvalid)
}

func Test_NoEndTime(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	fs := NewFlowStats(ctx)
	startTime := time.Now().Add(-300 * time.Second)
	fs.SendTimeMsg(StartTime, startTime)
	fs.SendUint64Msg(SendBytes, 100000000)

	// wait untill the receiver take all message
	for len(fs.toMsgReceiverChan) != 0 {
		time.Sleep(1 * time.Millisecond)
	}

	// assert bitrates
	_, sendBitrateStr := fs.SendBitrate()
	assert.Equal(t, "2.7 Mbps", sendBitrateStr)
}
