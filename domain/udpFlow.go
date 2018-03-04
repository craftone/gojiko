package domain

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain/stats"

	"github.com/craftone/gojiko/domain/ipemu"
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
	Arg        UdpEchoFlowArg
	session    *GtpSession
	toReceiver chan UDPpacket
	stats      stats.FlowStats
}

// sender is for goroutine
func (u *UdpEchoFlow) sender() {
	myLog := u.newMyLog("UdpFlowSender")
	myLog.Info("Start a UDP Flow sender goroutine")

	packetSize := u.Arg.SendPacketSize
	udpSize := packetSize - 20

	sendIntervalSec := float64(packetSize*8) / float64(u.Arg.TargetBps)
	sendInterval := time.Duration(sendIntervalSec * float64(time.Second))

	udpBody := make([]byte, udpSize)
	// 0 -  1 : Source Port
	// 2 -  3 : Destination Port
	// 4 -  5 : UDP length
	// 6 -  7 : checksum (ignore)
	binary.BigEndian.PutUint16(udpBody[0:], u.Arg.SourcePort)
	binary.BigEndian.PutUint16(udpBody[2:], uint16(u.Arg.DestAddr.Port))
	binary.BigEndian.PutUint16(udpBody[4:], udpSize)

	payload := udpBody[8:]
	// 0 -  1 : Receive Packet size (16bit)
	// 2 - 10 : Sequence Number (64bit)
	binary.BigEndian.PutUint16(payload[0:], u.Arg.RecvPacketSize)

	sess := u.session
	ipv4Emu := ipemu.NewIPv4Emulator(ipemu.UDP, sess.Paa(), u.Arg.DestAddr.IP, config.MTU())
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
				log.Info("Session status is not connected")
				break loop
			}
			seqNum++
			if seqNum > numOfSend {
				break loop
			}
			binary.BigEndian.PutUint64(payload[2:], seqNum)
			packet, err := ipv4Emu.NewIPv4GPDU(teid, u.Arg.Tos, u.Arg.Ttl, udpBody)
			if err != nil {
				myLog.Error(err)
			} else {
				myLog.Debugf("Send a packet #%d at %s", seqNum, time.Now())
				senderChan <- UDPpacket{sess.pgwDataAddr, packet}
			}
			nextTime = nextTime.Add(sendInterval)
			nextTimeChan = time.After(nextTime.Sub(time.Now()))
		}
	}
	sess.udpFlow = nil
	log.Info("End a UDP Flow goroutine")
}

// receiver is for goroutine
func (u *UdpEchoFlow) receiver() {
	myLog := u.newMyLog("UdpFlowReceiver")
	myLog.Info("Start a UDP Flow receiver goroutine")
	ipv4emu := ipemu.NewIPv4Emulator(ipemu.UDP, u.Arg.DestAddr.IP, u.session.Paa(), config.MTU())
	for pkt := range u.toReceiver {
		payload, err := ipv4emu.PickOutPayload(u.Arg.SourcePort, pkt.body)
		if err != nil {
			myLog.Debug(err)
			continue
		}
		seqNum := binary.BigEndian.Uint64(payload[2:])
		myLog.Debugf("Received #%d", seqNum)
	}
}

func (u *UdpEchoFlow) newMyLog(routine string) *logrus.Entry {
	sourceAddr := net.UDPAddr{IP: u.session.Paa(), Port: int(u.Arg.SourcePort)}
	return log.WithFields(logrus.Fields{
		"routine":        routine,
		"DestAddr":       u.Arg.DestAddr.String(),
		"SourceAddr":     sourceAddr.String(),
		"SendPacketSize": u.Arg.SendPacketSize,
		"TypeOfService":  u.Arg.Tos,
		"TTL":            u.Arg.Ttl,
		"TargetBps":      u.Arg.TargetBps,
		"NumOfSend":      u.Arg.NumOfSend,
		"RecvPacketSize": u.Arg.RecvPacketSize,
	})
}
