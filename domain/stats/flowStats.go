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
	return fs
}

// Duration returns duration in sec from StartTime to (EndTime | Now)
func (fs *FlowStats) Duration() float64 {
	endTime := fs.ReadTime(EndTime)
	if endTime.Equal(time.Time{}) {
		endTime = time.Now()
	}
	duration := endTime.Sub(fs.ReadTime(StartTime))
	return float64(duration) / float64(time.Second)
}

func (fs *FlowStats) bitrate(key Key) float64 {
	bitrate := float64(fs.ReadUint64(key)) * 8 / fs.Duration()
	return bitrate
}

func (fs *FlowStats) SendBitrate() (float64, string) {
	bitrate := fs.bitrate(SendBytes)
	str := FormatSIFloat(bitrate, "bps")
	return bitrate, str
}

func (fs *FlowStats) RecvBitrate() (float64, string) {
	bitrate := fs.bitrate(RecvBytes)
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
