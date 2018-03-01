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

type UdpEchoFlow struct {
	Arg UdpEchoFlowArg
}

func (u *UdpEchoFlow) sender(sess *GtpSession) {
	sourceAddr := net.UDPAddr{IP: sess.Paa(), Port: int(u.Arg.SourcePort)}
	myLog := log.WithFields(logrus.Fields{
		"routine":        "UdpFlowSender",
		"DestAddr":       u.Arg.DestAddr.String(),
		"SourceAddr":     sourceAddr.String(),
		"SendPacketSize": u.Arg.SendPacketSize,
		"TypeOfService":  u.Arg.Tos,
		"TTL":            u.Arg.Ttl,
		"TargetBps":      u.Arg.TargetBps,
		"NumOfSend":      u.Arg.NumOfSend,
		"RecvPacketSize": u.Arg.RecvPacketSize,
	})
	myLog.Debug("Start a UDP Flow goroutine")

	packetSize := u.Arg.SendPacketSize
	udpSize := packetSize - 20

	sendIntervalSec := float64(packetSize*8) / float64(u.Arg.TargetBps)
	sendInterval := time.Duration(sendIntervalSec * float64(time.Second))

	udpBody := make([]byte, udpSize)
	binary.BigEndian.PutUint16(udpBody[0:], u.Arg.SourcePort)
	binary.BigEndian.PutUint16(udpBody[2:], uint16(u.Arg.DestAddr.Port))
	binary.BigEndian.PutUint16(udpBody[4:], udpSize)
	binary.BigEndian.PutUint16(udpBody[8:], u.Arg.RecvPacketSize)

	ipv4Emu := ipemu.NewIPv4Emulator(ipemu.UDP, sourceAddr.IP, u.Arg.DestAddr.IP, 1500)
	log.Debugf("Make IPv4Emu : %#v", ipv4Emu)
	teid := sess.pgwDataFTEID.Teid()
	senderChan := sess.sgwCtrl.Pair().ToSender()
	seqNum := uint64(0)
	numOfSend := uint64(u.Arg.NumOfSend)

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
			packet, err := ipv4Emu.NewIPv4GPDU(teid, u.Arg.Tos, u.Arg.Ttl, udpBody)
			if err != nil {
				myLog.Debug(err)
			} else {
				myLog.Debugf("Send a packet #%d at %s", seqNum, time.Now())
				senderChan <- UDPpacket{sess.pgwDataAddr, packet}
			}
			nextTime = nextTime.Add(sendInterval)
			nextTimeChan = time.After(nextTime.Sub(time.Now()))
		}
	}
	sess.udpFlow = nil
	log.Debug("End a UDP Flow goroutine")
}
