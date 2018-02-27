package domain

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/craftone/gojiko/ipemu"
	"github.com/sirupsen/logrus"
)

const MIN_UDP_ECHO_PACKET_SIZE = 38

type UdpFlow struct {
	destAddr       net.UDPAddr
	sourcePort     uint16
	sendPacketSize uint16
	tos            byte
	ttl            byte
	targetBps      uint64
	sendDuration   time.Duration
	recvPacketSize uint16
}

func (u *UdpFlow) sender(sess *GtpSession) {
	sourceAddr := net.UDPAddr{IP: sess.paa.IPv4(), Port: int(u.sourcePort)}
	myLog := log.WithFields(logrus.Fields{
		"routine":        "UdpFlowSender",
		"DestAddr":       u.destAddr,
		"SourceAddr":     sourceAddr,
		"SendPacketSize": u.sendPacketSize,
		"TypeOfService":  u.tos,
		"TTL":            u.ttl,
		"TargetBps":      u.targetBps,
		"SendDuration":   u.sendDuration,
		"RecvPacketSize": u.recvPacketSize,
	})
	myLog.Debug("Start a UDP Flow goroutine")

	packetSize := u.sendPacketSize
	udpSize := packetSize - 20

	sendIntervalSec := float64(packetSize*8) / float64(u.targetBps)
	sendInterval := time.Duration(sendIntervalSec * float64(time.Second))

	udpBody := make([]byte, udpSize)
	binary.BigEndian.PutUint16(udpBody[0:], u.sourcePort)
	binary.BigEndian.PutUint16(udpBody[2:], uint16(u.destAddr.Port))
	binary.BigEndian.PutUint16(udpBody[4:], udpSize)
	binary.BigEndian.PutUint16(udpBody[8:], u.recvPacketSize)
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
	sess.udpFlow = nil
	log.Debug("End a UDP Flow goroutine")
}
