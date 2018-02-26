package domain

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/craftone/gojiko/ipemu"
	"github.com/sirupsen/logrus"
)

type UdpFlow struct {
	destAddr        net.UDPAddr
	sourcePort      uint16
	sendUdpDataSize uint16
	tos             byte
	ttl             byte
	targetBps       uint64
	sendDuration    time.Duration
	recvUdpDataSize uint16
}

func (u *UdpFlow) sender(sess *GtpSession) {
	sourceAddr := net.UDPAddr{IP: sess.paa.IPv4(), Port: int(u.sourcePort)}
	myLog := log.WithFields(logrus.Fields{
		"routine":         "UdpFlowSender",
		"DestAddr":        u.destAddr,
		"SourceAddr":      sourceAddr,
		"SendUdpDataSize": u.sendUdpDataSize,
		"TypeOfService":   u.tos,
		"TTL":             u.ttl,
		"TargetBps":       u.targetBps,
		"SendDuration":    u.sendDuration,
		"RecvUdpDataSize": u.recvUdpDataSize,
	})
	myLog.Debug("Start a UDP Flow goroutine")

	packetSize := u.sendUdpDataSize
	packetSize += 8  // UDP header
	packetSize += 20 // IP header

	sendIntervalSec := float64(packetSize*8) / float64(u.targetBps)
	sendInterval := time.Duration(sendIntervalSec * float64(time.Second))

	udpBody := make([]byte, u.sendUdpDataSize+8)
	binary.BigEndian.PutUint16(udpBody[0:], u.sourcePort)
	binary.BigEndian.PutUint16(udpBody[2:], uint16(u.destAddr.Port))
	binary.BigEndian.PutUint16(udpBody[4:], u.sendUdpDataSize+8)
	binary.BigEndian.PutUint16(udpBody[8:], u.recvUdpDataSize)
	ipv4Emu := ipemu.NewIPv4Emulator(ipemu.UDP, sess.Paa(), u.destAddr.IP, 1500)
	teid := sess.pgwDataFTEID.Teid()
	senderChan := sess.sgwCtrl.Pair().ToSender()
	seqNum := uint64(0)

	durationChan := time.After(u.sendDuration)

	nextTime := time.Now()
	nextTimeChan := time.After(0)
loop:
	select {
	case <-nextTimeChan:
		if sess.status != GssConnected {
			log.Debug("Session status is not connected")
			break
		}
		seqNum++
		binary.BigEndian.PutUint64(udpBody[10:], seqNum)
		packet, err := ipv4Emu.NewIPv4GPDU(teid, u.tos, u.ttl, udpBody)
		if err != nil {
			myLog.Debug(err)
		} else {
			myLog.Debugf("Send a packet at %s", time.Now())
			senderChan <- UDPpacket{sess.pgwDataAddr, packet}
		}
		nextTime = nextTime.Add(sendInterval)
		nextTimeChan = time.After(nextTime.Sub(time.Now()))
		goto loop
	case <-durationChan:
		log.Debug("UDP flow duration exipred")
	}
	log.Debug("End a UDP Flow goroutine")
}
