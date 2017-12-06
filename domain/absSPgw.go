package domain

import (
	"net"
	"sync"

	"github.com/craftone/gojiko/gtpv2c"

	"github.com/sirupsen/logrus"

	"github.com/craftone/gojiko/gtp"
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

	fromReceiver chan UDPpacket
	toSender     chan UDPpacket

	opSpgwMap map[string]*opSPgw //Key : UDPAddr.toString()
	mtxOp     sync.RWMutex
}

type SPgwIf interface {
	nextTeid() gtp.Teid
	nextSeqNum() uint32
	UDPAddr() net.UDPAddr
	getPair() SPgwIf
	findOrCreateOpSPgw(addr net.UDPAddr) (*opSPgw, error)
}

func newAbsSPgw(addr net.UDPAddr, recovery byte, pair SPgwIf) (*absSPgw, error) {
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	spgw := &absSPgw{
		addr:         addr,
		conn:         conn,
		recovery:     recovery,
		pair:         pair,
		fromReceiver: make(chan UDPpacket),
		toSender:     make(chan UDPpacket),
		opSpgwMap:    make(map[string]*opSPgw),
	}
	go absSPgwSenderRoutine(spgw, spgw.toSender)
	go absSPgwReceiverRoutine(spgw, spgw.fromReceiver)
	return spgw, nil
}

// absSPgwSenderRoutine is for GoRoutine
func absSPgwSenderRoutine(spgw *absSPgw, sendChan <-chan UDPpacket) {
	myLog := log.WithFields(logrus.Fields{
		"laddr":   spgw.addr,
		"routine": "SPgwSender",
	})
	myLog.Info("Start a SPgw Sender goroutine")

	for msg := range sendChan {
		myLog.Debug("Sending packet : ", msg)
		conn := spgw.conn
		_, err := conn.WriteToUDP(msg.body, &msg.raddr)
		if err != nil {
			myLog.Error(err)
			continue
		}
	}
}

// absSPgwReceiverRoutine is for GoRoutine
func absSPgwReceiverRoutine(spgw *absSPgw, recvChan chan<- UDPpacket) {
	myLog := log.WithFields(logrus.Fields{
		"laddr":   spgw.addr,
		"routine": "SPgwReceiver",
	})
	myLog.Info("Start a SPgw Receiver goroutine")

	buf := make([]byte, 2000)
	for {
		n, addr, err := spgw.conn.ReadFromUDP(buf)
		if err != nil {
			log.Error(err)
			continue
		}
		myLog.Debug("Received packet from %s : %v", buf[:n], addr.String())
		msg, _, err := gtpv2c.Unmarshal(buf[:n])
		if err != nil {
			log.Error(err)
			continue
		}
		myLog.Debug("Received GTPv2c MSG : %v", msg)
	}
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

func (sp *absSPgw) getPair() SPgwIf {
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
