package main

import (
	"fmt"
	"time"

	"github.com/craftone/gojiko/domain/stats"
	"github.com/dustin/go-humanize"
)

type sendRecvStatsSnapshot struct {
	sendPackets uint64
	sendBytes   uint64
	recvPackets uint64
	recvBytes   uint64
	timestamp   time.Time
}

func newSendRecvStatsSnapshot(s *sendRecvStats) sendRecvStatsSnapshot {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return sendRecvStatsSnapshot{
		sendBytes:   s.sendBytes,
		sendPackets: s.sendPackets,
		recvBytes:   s.recvBytes,
		recvPackets: s.recvPackets,
		timestamp:   time.Now(),
	}

}

func (s sendRecvStatsSnapshot) SubSendReport(l sendRecvStatsSnapshot) string {
	duration := s.timestamp.Sub(l.timestamp)
	durationSec := float64(duration) / float64(time.Second)

	sendBytes := s.sendBytes - l.sendBytes
	sendPackets := s.sendPackets - l.sendPackets

	return "TX" + s.report(durationSec, sendBytes, sendPackets)
}

func (s sendRecvStatsSnapshot) SubRecvReport(l sendRecvStatsSnapshot) string {
	duration := s.timestamp.Sub(l.timestamp)
	durationSec := float64(duration) / float64(time.Second)

	recvBytes := s.recvBytes - l.recvBytes
	recvPackets := s.recvPackets - l.recvPackets

	return "RX" + s.report(durationSec, recvBytes, recvPackets)
}

func (s sendRecvStatsSnapshot) report(durationSec float64, bytes, packets uint64) string {
	bps := float64(bytes*8) / durationSec
	pps := float64(packets) / durationSec
	return fmt.Sprintf(" in %.1f sec : %s / %s => %s / %s",
		durationSec,
		humanize.IBytes(bytes), stats.FormatSIUint(packets, "pkts"),
		stats.FormatSIFloat(bps, "bps"), stats.FormatSIFloat(pps, "pps"))
}
