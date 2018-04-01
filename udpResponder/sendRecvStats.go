package main

import (
	"fmt"
	"sync"

	"github.com/dustin/go-humanize"
)

type sendRecvStats struct {
	mtx         sync.RWMutex
	sendPackets uint64
	sendBytes   uint64
	recvPackets uint64
	recvBytes   uint64
}

var totalSendRecvStats = &sendRecvStats{}

func (s *sendRecvStats) String() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return fmt.Sprintf("RX: %s / %s pkts, TX: %s / %s pkts",
		humanize.IBytes(s.recvBytes), humanize.Comma(int64(s.recvPackets)),
		humanize.IBytes(s.sendBytes), humanize.Comma(int64(s.sendPackets)))
}

func (s *sendRecvStats) writeSend(packets, bytes uint64) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.sendPackets += packets
	s.sendBytes += bytes
}

func (s *sendRecvStats) writeRecv(packets, bytes uint64) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.recvPackets += packets
	s.recvBytes += bytes
}
