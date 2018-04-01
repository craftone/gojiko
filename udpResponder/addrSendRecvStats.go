package main

import (
	"fmt"
	"log"
	"sync"
)

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
func (a *addrSendRecvStats) Strings() []string {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	res := make([]string, 0, len(a.data))
	for key, val := range a.data {
		res = append(res, fmt.Sprintf("[%s] %s", key, val.String()))
	}
	return res
}

func (a *addrSendRecvStats) writeSend(addr string, packets, bytes uint64) {
	s := theAddrSendRecvStats.getSendRecvStats(addr)
	s.writeSend(packets, bytes)
}

func (a *addrSendRecvStats) writeRecv(addr string, packets, bytes uint64) {
	s := theAddrSendRecvStats.getSendRecvStats(addr)
	s.writeRecv(packets, bytes)
}
