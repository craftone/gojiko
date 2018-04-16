package domain

import (
	"net"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/craftone/gojiko/domain/gtp"
)

type UDPpacket struct {
	raddr net.UDPAddr
	body  []byte
}

type absSPgw struct {
	// Listen and Source UDP Address/Port
	addr       net.UDPAddr
	conn       *net.UDPConn
	recovery   byte
	teidVal    gtp.Teid
	seqNum     uint32
	mtxTeidSeq sync.Mutex
	pair       SPgwIf

	toSender       chan UDPpacket
	toEchoReceiver chan UDPpacket

	opSpgwMap map[string]*opSPgw //Key : UDPAddr.toString()
	mtxOp     sync.RWMutex
}

type SPgwIf interface {
	nextTeid() gtp.Teid
	nextSeqNum() uint32
	UDPAddr() net.UDPAddr
	Pair() SPgwIf
	Recovery() byte
	ToSender() chan UDPpacket
	findOrCreateOpSPgw(addr net.UDPAddr) (*opSPgw, error)
}

func newAbsSPgw(addr net.UDPAddr, recovery byte, pair SPgwIf) (*absSPgw, error) {
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	spgw := &absSPgw{
		addr:           addr,
		conn:           conn,
		recovery:       recovery,
		teidVal:        gtp.Teid(1),
		pair:           pair,
		toSender:       make(chan UDPpacket, 1000),
		toEchoReceiver: make(chan UDPpacket, 10),
		opSpgwMap:      make(map[string]*opSPgw),
	}
	go spgw.absSPgwSenderRoutine()
	return spgw, nil
}

// absSPgwSenderRoutine is for GoRoutine
func (sp *absSPgw) absSPgwSenderRoutine() {
	log := log.WithFields(logrus.Fields{
		"laddr":   sp.addr.String(),
		"routine": "SPgwSender",
	})
	log.Info("Start a SPgw Sender goroutine")

	for msg := range sp.toSender {
		log.WithField("raddr", msg.raddr.String()).Debugf("Sending %d bytes packet", len(msg.body))
		conn := sp.conn
		_, err := conn.WriteToUDP(msg.body, &msg.raddr)
		if err != nil {
			log.Error(err)
			continue
		}
	}
	log.Info("End a SPgw Sender goroutine")
}

func (sp *absSPgw) nextTeid() gtp.Teid {
	sp.mtxTeidSeq.Lock()
	defer sp.mtxTeidSeq.Unlock()
	teid := sp.teidVal
	sp.teidVal++
	return teid
}

func (sp *absSPgw) nextSeqNum() uint32 {
	sp.mtxTeidSeq.Lock()
	defer sp.mtxTeidSeq.Unlock()
	seqNum := sp.seqNum
	sp.seqNum++
	if sp.seqNum >= 0x800000 {
		sp.seqNum = 0
	}
	return seqNum
}

func (sp *absSPgw) UDPAddr() net.UDPAddr {
	return sp.addr
}

func (sp *absSPgw) Pair() SPgwIf {
	return sp.pair
}

func (sp *absSPgw) findOrCreateOpSPgw(addr net.UDPAddr) (*opSPgw, error) {
	sp.mtxOp.Lock()
	defer sp.mtxOp.Unlock()

	if val, ok := sp.opSpgwMap[addr.String()]; ok {
		return val, nil
	}
	opSPgw, err := newOpSPgw(sp, addr)
	if err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"laddr": opSPgw.parent.UDPAddr(),
		"raddr": addr.String(),
	}).Debug("Post an OpSPgw")
	sp.opSpgwMap[addr.String()] = opSPgw

	return opSPgw, nil
}

func (sp *absSPgw) ToSender() chan UDPpacket {
	return sp.toSender
}

func (sp *absSPgw) Recovery() byte {
	return sp.recovery
}
