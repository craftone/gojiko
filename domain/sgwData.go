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
		myLog.Debugf("Received packet from %s : %v", raddr.String(), buf[:n])

		if n < 8 {
			myLog.Errorf("Too short packet : %v", buf[:n])
			continue
		}
		msgType := gtpv1u.MessageTypeNum(buf[1])

		switch msgType {
		case gtpv1u.EchoRequestNum:
			myLog.Error("Not yet implemented!")
			// sgwCtrl.toEchoReceiver <- UDPpacket{*raddr, buf[:n]}
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
			copy(received, buf[:n])
			sess.fromSgwDataReceiverChan <- UDPpacket{*raddr, received}
		default:
			myLog.Debugf("Unkown Message Type : %d", msgType)
		}
	}
}
