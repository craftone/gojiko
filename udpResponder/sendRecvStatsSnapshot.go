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

func (s sendRecvStatsSnapshot) reportBytesPkts(l sendRecvStatsSnapshot) string {
	duration := s.timestamp.Sub(l.timestamp)
	durationSec := float64(duration) / float64(time.Second)

	sendBytes := s.sendBytes - l.sendBytes
	sendPackets := s.sendPackets - l.sendPackets
	recvBytes := s.recvBytes - l.recvBytes
	recvPackets := s.recvPackets - l.recvPackets

	return fmt.Sprintf("RX=>TX Bytes/pkts in %.1f sec : %s / %s => %s / %s",
		durationSec,
		humanize.IBytes(recvBytes), stats.FormatSIUint(recvPackets, "pkts"),
		humanize.IBytes(sendBytes), stats.FormatSIUint(sendPackets, "pkts"))
}

func (s sendRecvStatsSnapshot) reportBpsPps(l sendRecvStatsSnapshot) string {
	duration := s.timestamp.Sub(l.timestamp)
	durationSec := float64(duration) / float64(time.Second)

	sendBytes := s.sendBytes - l.sendBytes
	sendBps := float64(sendBytes*8) / durationSec
	sendPackets := s.sendPackets - l.sendPackets
	sendPps := float64(sendPackets) / durationSec
	recvBytes := s.recvBytes - l.recvBytes
	recvBps := float64(recvBytes*8) / durationSec
	recvPackets := s.recvPackets - l.recvPackets
	recvPps := float64(recvPackets) / durationSec

	return fmt.Sprintf("RX=>TX bps/pps    in %.1f sec : %s / %s => %s / %s",
		durationSec,
		stats.FormatSIFloat(recvBps, "bps"), stats.FormatSIFloat(recvPps, "pps"),
		stats.FormatSIFloat(sendBps, "bps"), stats.FormatSIFloat(sendPps, "pps"))
}
