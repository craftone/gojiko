package domain

import (
	"net"
	"sync"

	"github.com/craftone/gojiko/gtp"
)

type absSPgw struct {
	// Listen and Source UDP Address/Port
	addr     net.UDPAddr
	recovery byte
	teidVal  gtp.Teid
	mtx      sync.Mutex

	fromRceiver <-chan []byte
	toSender    chan<- []byte
	opSpgwMap   map[string]SPgw //Key : UDPAddr.toString()
}

type SPgw interface {
	nextTeid() gtp.Teid
}

func newAbsSPgw(addr net.UDPAddr, recovery byte) *absSPgw {
	return &absSPgw{
		addr:        addr,
		recovery:    recovery,
		fromRceiver: make(chan []byte),
		toSender:    make(chan []byte),
	}
}

func (sp *absSPgw) nextTeid() gtp.Teid {
	sp.mtx.Lock()
	defer sp.mtx.Unlock()
	teid := sp.teidVal
	sp.teidVal++
	return teid
}
