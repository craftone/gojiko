package main

import (
	"fmt"
	"log"
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

type addrSendRecvStats struct {
	mtx  sync.RWMutex
	data map[string]*sendRecvStats
}

var theAddrSendRecvStats = &addrSendRecvStats{data: make(map[string]*sendRecvStats)}

func (a *addrSendRecvStats) getSendRecvStats(addr string) *sendRecvStats {
	if addr == "" {
		log.Panic("blank addr")
	}
	a.mtx.Lock()
	defer a.mtx.Unlock()
	if s, ok := a.data[addr]; ok {
		return s
	}
	s := &sendRecvStats{}
	a.data[addr] = s
	return s
}

func (s *sendRecvStats) String() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return fmt.Sprintf("RX: %s / %s pkts, TX: %s / %s pkts",
		humanize.IBytes(s.recvBytes), humanize.Comma(int64(s.recvPackets)),
		humanize.IBytes(s.sendBytes), humanize.Comma(int64(s.sendPackets)))
}

func (a *addrSendRecvStats) Strings() []string {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	res := make([]string, 0, len(a.data))
	for key, val := range a.data {
		res = append(res, fmt.Sprintf("[%s] %s", key, val.String()))
	}
	return res
}

func writeSend(addr string, packets, bytes uint64) {
	s := theAddrSendRecvStats.getSendRecvStats(addr)
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.sendPackets += packets
	s.sendBytes += bytes
}

func writeRecv(addr string, packets, bytes uint64) {
	s := theAddrSendRecvStats.getSendRecvStats(addr)
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.recvPackets += packets
	s.recvBytes += bytes
}

// func read(addr string) sendRecvStats {
// 	s := theAddrSendRecvStats.getSendRecvStats(addr)
// 	return sendRecvStats{
// 		sendPackets: s.sendPackets,
// 		sendBytes:   s.sendBytes,
// 		recvPackets: s.recvPackets,
// 		recvBytes:   s.recvBytes,
// 	}
// }
