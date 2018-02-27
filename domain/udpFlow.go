package domain

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/craftone/gojiko/ipemu"
	"github.com/sirupsen/logrus"
)

const MIN_UDP_ECHO_PACKET_SIZE = 38

type UdpEchoFlowArg struct {
	DestAddr       net.UDPAddr
	SourcePort     uint16
	SendPacketSize uint16
	Tos            byte
	Ttl            byte
	TargetBps      uint64
	NumOfSend      int
	RecvPacketSize uint16
}

func (u *UdpEchoFlowArg) sender(sess *GtpSession) {
	sourceAddr := net.UDPAddr{IP: sess.paa.IPv4(), Port: int(u.SourcePort)}
	myLog := log.WithFields(logrus.Fields{
		"routine":        "UdpFlowSender",
		"DestAddr":       u.DestAddr,
		"SourceAddr":     sourceAddr,
		"SendPacketSize": u.SendPacketSize,
		"TypeOfService":  u.Tos,
		"TTL":            u.Ttl,
		"TargetBps":      u.TargetBps,
		"NumOfSend":      u.NumOfSend,
		"RecvPacketSize": u.RecvPacketSize,
	})
	myLog.Debug("Start a UDP Flow goroutine")

	packetSize := u.SendPacketSize
	udpSize := packetSize - 20

	sendIntervalSec := float64(packetSize*8) / float64(u.TargetBps)
	sendInterval := time.Duration(sendIntervalSec * float64(time.Second))

	udpBody := make([]byte, udpSize)
	binary.BigEndian.PutUint16(udpBody[0:], u.SourcePort)
	binary.BigEndian.PutUint16(udpBody[2:], uint16(u.DestAddr.Port))
	binary.BigEndian.PutUint16(udpBody[4:], udpSize)
	binary.BigEndian.PutUint16(udpBody[8:], u.RecvPacketSize)
	ipv4Emu := ipemu.NewIPv4Emulator(ipemu.UDP, sess.Paa(), u.DestAddr.IP, 1500)
	teid := sess.pgwDataFTEID.Teid()
	senderChan := sess.sgwCtrl.Pair().ToSender()
	seqNum := uint64(0)
	numOfSend := uint64(u.NumOfSend)

	nextTime := time.Now()
	nextTimeChan := time.After(0)

loop:
	for {
		select {
		case <-nextTimeChan:
			if sess.status != GssConnected {
				log.Debug("Session status is not connected")
				break loop
			}
			seqNum++
			if seqNum > numOfSend {
				break loop
			}
			binary.BigEndian.PutUint64(udpBody[10:], seqNum)
			packet, err := ipv4Emu.NewIPv4GPDU(teid, u.Tos, u.Ttl, udpBody)
			if err != nil {
				myLog.Debug(err)
			} else {
				myLog.Debugf("Send a packet at %s", time.Now())
				senderChan <- UDPpacket{sess.pgwDataAddr, packet}
			}
			nextTime = nextTime.Add(sendInterval)
			nextTimeChan = time.After(nextTime.Sub(time.Now()))
		}
	}
	sess.udpFlow = nil
	log.Debug("End a UDP Flow goroutine")
}
