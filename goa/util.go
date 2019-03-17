package main

import (
	"fmt"
	"net"
	"time"

	"github.com/craftone/gojiko/domain/stats"

	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/goa/app"
)

func newFteid(ip net.IP, teid gtp.Teid) *app.Fteid {
	return &app.Fteid{Ipv4: ip.String(), Teid: fmt.Sprintf("0x%08X", teid)}
}

func querySgw(sgwAddr string) (*domain.SgwCtrl, error) {
	sgwCtrlAddr := net.UDPAddr{IP: net.ParseIP(sgwAddr), Port: domain.GtpControlPort}
	theSgwCtrlRepo := domain.TheSgwCtrlRepo()
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(sgwCtrlAddr)
	if sgwCtrl == nil {
		return nil, fmt.Errorf("There is no SGW whose IP address is %s", sgwCtrlAddr.String())
	}
	return sgwCtrl, nil
}

func querySessionByIMSIandEBI(sgwAddr, imsi string, ebi int) (*domain.GtpSession, error) {
	sgwCtrl, err := querySgw(sgwAddr)
	if err != nil {
		return nil, err
	}
	sess := sgwCtrl.FindByImsiEbi(imsi, byte(ebi))
	if sess == nil {
		return nil, fmt.Errorf("There is no session whose IMSI is %s and EBI is %d", imsi, ebi)
	}
	return sess, nil
}

func newGtpsessionMedia(sess *domain.GtpSession) *app.Gtpsession {
	return &app.Gtpsession{
		Apn: sess.Apn(),
		Ebi: int(sess.Ebi()),
		Fteid: &app.GtpSessionFTEIDs{
			PgwCtrlFTEID: newFteid(sess.PgwCtrlFTEID()),
			PgwDataFTEID: newFteid(sess.PgwDataFTEID()),
			SgwCtrlFTEID: newFteid(sess.SgwCtrlFTEID()),
			SgwDataFTEID: newFteid(sess.SgwDataFTEID()),
		},
		Paa:    sess.Paa().String(),
		Sid:    int(sess.ID()),
		Imsi:   sess.Imsi(),
		Mcc:    sess.Mcc(),
		Mei:    sess.Mei(),
		Mnc:    sess.Mnc(),
		Msisdn: sess.Msisdn(),
		Tai: &app.Tai{
			Mcc: sess.TaiMcc(),
			Mnc: sess.TaiMnc(),
			Tac: int(sess.TaiTac()),
		},
		Ecgi: &app.Ecgi{
			Mcc: sess.EcgiMcc(),
			Mnc: sess.EcgiMnc(),
			Eci: int(sess.EcgiEci()),
		},
		RatType: &app.RatType{
			RatTypeValue: int(sess.RatTypeValue()),
			RatType:      sess.RatType(),
		},
	}
}

func newUDPEchoFlowPayload(arg domain.UdpEchoFlowArg) *app.UDPEchoFlowPayload {
	return &app.UDPEchoFlowPayload{
		DestAddr:       arg.DestAddr.IP.String(),
		DestPort:       arg.DestAddr.Port,
		NumOfSend:      arg.NumOfSend,
		RecvPacketSize: int(arg.RecvPacketSize),
		SendPacketSize: int(arg.SendPacketSize),
		SourcePort:     int(arg.SourcePort),
		TargetBps:      int(arg.TargetBps),
		Tos:            int(arg.Tos),
		TTL:            int(arg.Ttl),
	}
}

func newStatsMedia(sts *stats.FlowStats) *app.SendRecvStatistics {
	sendBitrate, sendBitrateStr := sts.SendBitrate()
	recvBitrate, recvBitrateStr := sts.RecvBitrate()
	sendBytes, sendBytesStr := sts.SendBytes()
	recvBytes, recvBytesStr := sts.RecvBytes()
	sendPackets, sendPacketsStr := sts.SendPackets()
	recvPackets, recvPacketsStr := sts.RecvPackets()
	sendBytesSkipped, sendBytesSkippedStr := sts.SendBytesSkipped()
	recvBytesInvalid, recvBytesInvalidStr := sts.RecvBytesInvalid()
	sendPktsSkipped, sendPktsSkippedStr := sts.SendPacketsSkipped()
	recvPktsInvalid, recvPktsInvalidStr := sts.RecvPacketsInvalid()

	statsMedia := &app.SendRecvStatistics{
		StartTime: sts.ReadTime(stats.StartTime),
		Duration:  sts.Duration(),
		SendStats: &app.SendStatForMachie{
			Bitrate:        sendBitrate,
			Bytes:          int(sendBytes),
			Packets:        int(sendPackets),
			SkippedBytes:   int(sendBytesSkipped),
			SkippedPackets: int(sendPktsSkipped),
		},
		SendStatsHumanize: &app.SendStatForHuman{
			Bitrate:        sendBitrateStr,
			Bytes:          sendBytesStr,
			Packets:        sendPacketsStr,
			SkippedBytes:   sendBytesSkippedStr,
			SkippedPackets: sendPktsSkippedStr,
		},
		RecvStats: &app.RecvStatForMachie{
			Bitrate:        recvBitrate,
			Bytes:          int(recvBytes),
			Packets:        int(recvPackets),
			InvalidBytes:   int(recvBytesInvalid),
			InvalidPackets: int(recvPktsInvalid),
		},
		RecvStatsHumanize: &app.RecvStatForHuman{
			Bitrate:        recvBitrateStr,
			Bytes:          recvBytesStr,
			Packets:        recvPacketsStr,
			InvalidBytes:   recvBytesInvalidStr,
			InvalidPackets: recvPktsInvalidStr,
		},
	}
	if !sts.ReadTime(stats.EndTime).Equal(time.Time{}) {
		endTime := sts.ReadTime(stats.EndTime)
		statsMedia.EndTime = &endTime
	}
	return statsMedia
}

func newCauseMedia(gsRes domain.GsRes) *app.Gtpv2cCause {
	return &app.Gtpv2cCause{
		Type:   gsRes.Code.String(),
		Value:  int(gsRes.Value),
		Detail: gsRes.Msg,
	}
}
