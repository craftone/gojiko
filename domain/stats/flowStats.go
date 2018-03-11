package stats

import (
	"context"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/craftone/gojiko/config"
)

type FlowStats struct {
	*absStats
}

func NewFlowStats(ctx context.Context) *FlowStats {
	fs := &FlowStats{newAbsStats(ctx, config.StatsFlowBufferLength())}
	fs.SendTimeMsg(StartTime, time.Now())
	return fs
}

func (fs *FlowStats) bitrate(key Key, endTime time.Time) float64 {
	duration := endTime.Sub(fs.ReadTime(StartTime))
	durationSec := float64(duration) / float64(time.Second)
	bitrate := float64(fs.ReadUint64(key)) * 8 / durationSec
	return bitrate
}

func (fs *FlowStats) SendBitrate(endTime time.Time) (float64, string) {
	bitrate := fs.bitrate(SendBytes, endTime)
	str := FormatSIFloat(bitrate, "bps")
	return bitrate, str
}

func (fs *FlowStats) RecvBitrate(endTime time.Time) (float64, string) {
	bitrate := fs.bitrate(RecvBytes, endTime)
	str := FormatSIFloat(bitrate, "bps")
	return bitrate, str
}

func (fs *FlowStats) SendBytes() (uint64, string) {
	b := uint64(fs.ReadUint64(SendBytes))
	return b, humanize.IBytes(b)
}

func (fs *FlowStats) RecvBytes() (uint64, string) {
	b := uint64(fs.ReadUint64(RecvBytes))
	return b, humanize.IBytes(b)
}

func (fs *FlowStats) SendPackets() (uint64, string) {
	p := uint64(fs.ReadUint64(SendPackets))
	return p, FormatSIUint(p, "pkts")
}

func (fs *FlowStats) RecvPackets() (uint64, string) {
	p := uint64(fs.ReadUint64(RecvPackets))
	return p, FormatSIUint(p, "pkts")
}

func (fs *FlowStats) SendBytesSkipped() (uint64, string) {
	b := uint64(fs.ReadUint64(SendBytesSkipped))
	return b, humanize.IBytes(b)
}

func (fs *FlowStats) RecvBytesInvalid() (uint64, string) {
	b := uint64(fs.ReadUint64(RecvBytesInvalid))
	return b, humanize.IBytes(b)
}

func (fs *FlowStats) SendPacketsSkipped() (uint64, string) {
	p := uint64(fs.ReadUint64(SendPacketsSkipped))
	return p, FormatSIUint(p, "pkts")
}

func (fs *FlowStats) RecvPacketsInvalid() (uint64, string) {
	p := uint64(fs.ReadUint64(RecvPacketsInvalid))
	return p, FormatSIUint(p, "pkts")
}
