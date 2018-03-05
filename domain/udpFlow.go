package domain

import (
	"context"
	"encoding/binary"
	"net"
	"time"

	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain/stats"
	humanize "github.com/dustin/go-humanize"

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
	Arg                  UdpEchoFlowArg
	session              *GtpSession
	fromSessDataReceiver chan UDPpacket
	stats                *stats.FlowStats
	statsCtxCencel       context.CancelFunc
}

// sender is for goroutine
func (u *UdpEchoFlow) sender(ctx context.Context) {
	log := u.newMyLog("UdpFlowSender")
	log.Info("Start a UDP Flow sender goroutine")

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
	startTime := nextTime
	nextTimeChan := time.After(0)
	skipFlg := false // This flag becomes true when skipping due to delay etc.

loop:
	for {
		select {
		case <-nextTimeChan:
			if sess.status != GssConnected {
				log.Info("Session status is not connected")
				sess.udpFlow.statsCtxCencel()
				break loop
			}
			seqNum++
			if seqNum > numOfSend {
				time.Sleep(config.FlowUdpEchoWaitReceive())
				sess.udpFlow.statsCtxCencel()
				break loop
			}
			if !skipFlg {
				binary.BigEndian.PutUint64(payload[2:], seqNum)
				packet, err := ipv4Emu.NewIPv4GPDU(teid, u.Arg.Tos, u.Arg.Ttl, udpBody)
				if err != nil {
					log.Error(err)
				} else {
					senderChan <- UDPpacket{sess.pgwDataAddr, packet}
					log.Debugf("Send a packet #%d at %s", seqNum, time.Now())
					u.stats.SendInt64Msg(stats.SendPackets, 1)
					u.stats.SendInt64Msg(stats.SendBytes, int64(28+packetSize))
				}
			} else {
				u.stats.SendInt64Msg(stats.SendPacketsSkipped, 1)
				u.stats.SendInt64Msg(stats.SendBytesSkipped, int64(28+packetSize))
			}

			nextTime = nextTime.Add(sendInterval)
			if nextTime.Before(time.Now()) {
				skipFlg = true
				nextTimeChan = time.After(0)
			} else {
				skipFlg = false
				nextTimeChan = time.After(nextTime.Sub(time.Now()))
			}
		case <-ctx.Done():
			break loop
		}
	}
	endTime := time.Now()

	log.Info("End a UDP Flow sender goroutine")
	durationSec := float64(endTime.Sub(startTime)) / float64(time.Second)
	sendBytes := u.stats.ReadInt64(stats.SendBytes)
	log.Infof("[SEND stats] %s in %s(s) : %s / %s : (skipped) %s / %s",
		humanize.SI(float64(sendBytes)*8/durationSec, "bps"),
		humanize.Ftoa(durationSec),
		humanize.SI(float64(u.stats.ReadInt64(stats.SendBytes)), "bytes"),
		humanize.SI(float64(u.stats.ReadInt64(stats.SendPackets)), "pkts"),
		humanize.SI(float64(u.stats.ReadInt64(stats.SendBytesSkipped)), "bytes"),
		humanize.SI(float64(u.stats.ReadInt64(stats.SendPacketsSkipped)), "pkts"))
	sess.udpFlow = nil
}

// receiver is for goroutine
func (u *UdpEchoFlow) receiver(ctx context.Context) {
	myLog := u.newMyLog("UdpFlowReceiver")
	myLog.Info("Start a UDP Flow receiver goroutine")
	ipv4emu := ipemu.NewIPv4Emulator(ipemu.UDP, u.Arg.DestAddr.IP, u.session.Paa(), config.MTU())
	startTime := time.Now()
loop:
	select {
	case pkt := <-u.fromSessDataReceiver:
		payload, err := ipv4emu.PickOutPayload(u.Arg.SourcePort, pkt.body)
		if err != nil {
			u.stats.SendInt64Msg(stats.RecvPacketsInvalid, 1)
			u.stats.SendInt64Msg(stats.RecvBytesInvalid, int64(20+len(pkt.body)))
			myLog.Debug(err)
			goto loop
		}
		seqNum := binary.BigEndian.Uint64(payload[2:])
		u.stats.SendInt64Msg(stats.RecvPackets, 1)
		u.stats.SendInt64Msg(stats.RecvBytes, int64(20+len(pkt.body)))
		myLog.Debugf("Received #%d", seqNum)
		goto loop
	case <-ctx.Done():
		break loop
	}
	endTime := time.Now()
	log.Info("End a UDP Flow receiver goroutine")
	durationSec := float64(endTime.Sub(startTime)) / float64(time.Second)
	recvBytes := u.stats.ReadInt64(stats.RecvBytes)
	log.Infof("[RECV stats] %s in %s(s) : %s / %s : (invalid) %s / %s",
		humanize.SI(float64(recvBytes)*8/durationSec, "bps"),
		humanize.Ftoa(durationSec),
		humanize.SI(float64(u.stats.ReadInt64(stats.RecvBytes)), "bytes"),
		humanize.SI(float64(u.stats.ReadInt64(stats.RecvPackets)), "pkts"),
		humanize.SI(float64(u.stats.ReadInt64(stats.RecvBytesInvalid)), "bytes"),
		humanize.SI(float64(u.stats.ReadInt64(stats.RecvPacketsInvalid)), "pkts"))
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
