package domain

import (
	"context"
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
	Arg                  UdpEchoFlowArg
	session              *GtpSession
	fromSessDataReceiver chan UDPpacket
	stats                *stats.FlowStats
	ctxCencel            context.CancelFunc
}

// sender is for goroutine
func (u *UdpEchoFlow) sender(ctx context.Context) {
	log := u.newMyLog("UdpFlowSender")
	log.Info("Start a UDP Flow sender goroutine")

	packetSize := u.Arg.SendPacketSize
	udpSize := packetSize - 20

	sendIntervalSec := float64(packetSize*8) / float64(u.Arg.TargetBps)
	sendInterval := time.Duration(sendIntervalSec * float64(time.Second))

	sess := u.session
	ipv4Emu := ipemu.NewIPv4Emulator(ipemu.UDP, sess.Paa(), u.Arg.DestAddr.IP, config.MTU())
	teid := sess.pgwDataFTEID.Teid()
	senderChan := sess.sgwCtrl.Pair().ToSender()
	seqNum := uint64(0)
	numOfSend := uint64(u.Arg.NumOfSend)

	// pseudo IPv4 packet to calculate UDP checksum
	pseudoIPv4 := make([]byte, udpSize+12)
	// make ipv4 pseudo header
	copy(pseudoIPv4[0:4], sess.Paa())              // Source IPv4 Address
	copy(pseudoIPv4[4:8], u.Arg.DestAddr.IP.To4()) // Destination IPv4 Address
	pseudoIPv4[9] = byte(ipemu.UDP)                // Protocol : UDP(17)
	binary.BigEndian.PutUint16(                    // UDP length
		pseudoIPv4[10:12], uint16(udpSize))

	udpBody := pseudoIPv4[12:]
	// 0 -  1 : Source Port
	binary.BigEndian.PutUint16(udpBody[0:], u.Arg.SourcePort)
	// 2 -  3 : Destination Port
	binary.BigEndian.PutUint16(udpBody[2:], uint16(u.Arg.DestAddr.Port))
	// 4 -  5 : UDP length
	binary.BigEndian.PutUint16(udpBody[4:], udpSize)
	// 6 -  7 : checksum (ignore)
	payload := udpBody[8:]
	// 0 -  1 : Receive Packet size (16bit)
	binary.BigEndian.PutUint16(payload[0:], u.Arg.RecvPacketSize)
	// 2 - 10 : Sequence Number (64bit)
	//   set later

	// 12: Pseudo IPv4 header,
	//  8: UDP header
	//  2: RecvPacketSize
	// BIT-NOT(^) is nessesary because checksum did BIT-NOT at last but
	// udpChecksumBase should not be BIT-NOTed.
	udpChecksumBase := ^ipemu.Checksum(0, pseudoIPv4[0:12+8+2])

	nextTime := time.Now()
	nextTimeChan := time.After(0)

	// skipFlg represents that sending a packet will be skipped at this time of loop
	// due to delay etc, but count up seqNum etc, should be processed.
	skipFlg := false

loop:
	for {
		select {
		case <-nextTimeChan:
			if sess.status != GssConnected {
				log.Info("The session is not connecting")
				sess.StopUDPFlow()
				break loop
			}
			seqNum++
			if seqNum > numOfSend {
				time.Sleep(config.FlowUdpEchoWaitReceive())
				sess.StopUDPFlow()
				break loop
			}
			if !skipFlg {
				binary.BigEndian.PutUint64(payload[2:], seqNum)
				// calc udp checksum
				udpChecksum := ipemu.Checksum(udpChecksumBase, payload[2:10])
				if udpChecksum == 0 {
					udpChecksum = 0xFFFF
				}
				binary.BigEndian.PutUint16(udpBody[6:8], udpChecksum)

				packet, err := ipv4Emu.NewIPv4GPDU(teid, u.Arg.Tos, u.Arg.Ttl, udpBody)
				if err != nil {
					log.Error(err)
				} else {
					senderChan <- UDPpacket{sess.pgwDataAddr, packet}
					log.Debugf("Send a packet #%d at %s", seqNum, time.Now())
					u.stats.SendUint64Msg(stats.SendPackets, 1)
					u.stats.SendUint64Msg(stats.SendBytes, 28+uint64(packetSize))
				}
			} else {
				u.stats.SendUint64Msg(stats.SendPacketsSkipped, 1)
				u.stats.SendUint64Msg(stats.SendBytesSkipped, 28+uint64(packetSize))
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

	log.Info("End a UDP Flow sender goroutine")
	durationSec := u.stats.Duration()
	_, bitrate := u.stats.SendBitrate()
	_, sendBytes := u.stats.SendBytes()
	_, sendPackets := u.stats.SendPackets()
	_, sendBytesSkipped := u.stats.SendBytesSkipped()
	_, sendPacketsSkipped := u.stats.SendPacketsSkipped()
	log.Infof("[SEND stats] %s in %1.1f(s) : %s / %s : (skipped) %s / %s",
		bitrate, durationSec,
		sendBytes, sendPackets,
		sendBytesSkipped, sendPacketsSkipped)
}

// receiver is for goroutine
func (u *UdpEchoFlow) receiver(ctx context.Context) {
	myLog := u.newMyLog("UdpFlowReceiver")
	myLog.Info("Start a UDP Flow receiver goroutine")
	ipv4emu := ipemu.NewIPv4Emulator(ipemu.UDP, u.Arg.DestAddr.IP, u.session.Paa(), config.MTU())
loop:
	select {
	case pkt := <-u.fromSessDataReceiver:
		payload, err := ipv4emu.PickOutPayload(u.Arg.SourcePort, pkt.body)
		if err != nil {
			u.stats.SendUint64Msg(stats.RecvPacketsInvalid, 1)
			u.stats.SendUint64Msg(stats.RecvBytesInvalid, 20+uint64(len(pkt.body)))
			myLog.Debug(err)
			goto loop
		}
		seqNum := binary.BigEndian.Uint64(payload[2:])
		u.stats.SendUint64Msg(stats.RecvPackets, 1)
		u.stats.SendUint64Msg(stats.RecvBytes, 20+uint64(len(pkt.body)))
		myLog.Debugf("Received #%d", seqNum)
		goto loop
	case <-ctx.Done():
		break loop
	}
	log.Info("End a UDP Flow receiver goroutine")
	durationSec := u.stats.Duration()
	_, bitrate := u.stats.RecvBitrate()
	_, recvBytes := u.stats.RecvBytes()
	_, recvPackets := u.stats.RecvPackets()
	_, recvBytesInvalid := u.stats.RecvBytesInvalid()
	_, recvPacketsInvalid := u.stats.RecvPacketsInvalid()
	log.Infof("[RECV stats] %s in %1.1f(s) : %s / %s : (invalid) %s / %s",
		bitrate, durationSec,
		recvBytes, recvPackets,
		recvBytesInvalid, recvPacketsInvalid)
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

// getters and setters

func (u *UdpEchoFlow) Stats() *stats.FlowStats {
	return u.stats
}
