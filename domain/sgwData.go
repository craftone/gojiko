package domain

import (
	"encoding/binary"
	"net"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv1u"
	"github.com/sirupsen/logrus"
)

type SgwData struct {
	*absSPgw
}

func newSgwData(addr net.UDPAddr, recovery byte, sgwCtrl *SgwCtrl) (*SgwData, error) {
	myLog := log.WithFields(logrus.Fields{
		"addr":     addr.String(),
		"recovery": recovery,
	})
	absSPgw, err := newAbsSPgw(addr, recovery, sgwCtrl)
	if err != nil {
		return nil, err
	}
	myLog.Info("A new SGW Data is created")

	sgwData := &SgwData{absSPgw}
	go sgwData.sgwDataReceiverRoutine()
	go sgwData.echoReceiver()

	return sgwData, nil
}

// sgwDataReceiverRoutine is for GoRoutine
func (sgwData *SgwData) sgwDataReceiverRoutine() {
	myLog := log.WithFields(logrus.Fields{
		"laddr":   sgwData.addr.String(),
		"routine": "SgwDataReceiver",
	})
	myLog.Info("Start a SGW Data Receiver goroutine")

	buf := make([]byte, 2000)
	for {
		n, raddr, err := sgwData.conn.ReadFromUDP(buf)
		if err != nil {
			myLog.Error(err)
			continue
		}
		myLog.Debugf("Received a %d bytes packet from %s", n, raddr.String())

		if n < 8 {
			myLog.Errorf("Too short packet : %v", buf[:n])
			continue
		}
		msgType := gtpv1u.MessageTypeNum(buf[1])

		switch msgType {
		case gtpv1u.EchoRequestNum:
			myLog.Error("Not yet implemented!")
			sgwData.toEchoReceiver <- UDPpacket{*raddr, buf[:n]}
		case gtpv1u.EchoResponseNum:
			myLog.Error("Not yet implemented!")
			// Not yet be implemented
		case gtpv1u.GpduNum:
			teid := gtp.Teid(binary.BigEndian.Uint32(buf[4:8]))
			sess := sgwData.Pair().(*SgwCtrl).FindByDataTeid(teid)
			if sess == nil {
				myLog.Debug("No session that have the data teid : %04x", teid)
				continue
			}
			received := make([]byte, n)
			copy(received, buf[8:n])
			sess.fromSgwDataReceiverChan <- UDPpacket{*raddr, received}
		default:
			myLog.Debugf("Unkown Message Type : %d", msgType)
		}
	}
}

// echoReceiver is for GoRoutine
func (sd *SgwData) echoReceiver() {
	myLog := log.WithFields(logrus.Fields{
		"laddr":   sd.addr.String(),
		"routine": "SgwDataEchoReceiver",
	})
	myLog.Info("Start a SPgw ECHO Receiver goroutine")

	for pkt := range sd.toEchoReceiver {
		// Checking valid GTPv1-U is ommited since it is not important.
		myLog.Debugf("Received ECHO Request from %#v", pkt.raddr)

		seqNum := binary.BigEndian.Uint16(pkt.body[8:10])
		// make ECHO Response
		body := []byte{
			byte(1<<5 + 1<<4 + 1<<1), // Version:1, PT:1, Sequence:1
			byte(gtpv1u.EchoResponseNum),
			0, 6, // Length
			0, 0, 0, 0, // TEID is always 0
			byte(seqNum >> 8), byte(seqNum), // Sequence Number
			0,           // N-PDU Number
			0,           // Next Extention Header Type
			14,          // RecoveryNum
			sd.recovery, // Recovery Value
		}

		res := UDPpacket{
			raddr: pkt.raddr,
			body:  body,
		}
		sd.toSender <- res
	}
}
