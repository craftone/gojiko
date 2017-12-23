package domain

import (
	"net"
	"sync"
	"time"

	gsc "github.com/craftone/gojiko/domain/gtpSessionCmd"
	"github.com/sirupsen/logrus"

	"github.com/craftone/gojiko/gtpv2c"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

type gtpSessionStatus byte

const (
	gssIdle gtpSessionStatus = iota
	gssCSReqSending
	gssConnected
)

type gtpSession struct {
	id     SessionID
	status gtpSessionStatus
	mtx    sync.RWMutex

	cmdReqChan           chan gsc.Cmd
	cmdResChan           chan gsc.Res
	toCtrlSenderChan     chan UDPpacket
	fromCtrlReceiverChan chan UDPpacket
	toDataSenderChan     chan UDPpacket
	fromDataReceiverChan chan UDPpacket

	sgwCtrl      *SgwCtrl
	sgwCtrlFTEID *ie.Fteid
	sgwDataFTEID *ie.Fteid
	pgwCtrlFTEID *ie.Fteid // at first, ZERO FTEID
	pgwDataFTEID *ie.Fteid // at first, nil

	pgwCtrlAddr net.UDPAddr
	pgwDataAddr net.UDPAddr
	sgwCtrlAddr net.UDPAddr
	sgwDataAddr net.UDPAddr

	imsi           *ie.Imsi
	msisdn         *ie.Msisdn
	ebi            *ie.Ebi
	paa            *ie.Paa // at first, 0.0.0.0
	apn            *ie.Apn
	ambr           *ie.Ambr
	ratType        *ie.RatType
	servingNetwork *ie.ServingNetwork
	pdnType        *ie.PdnType
}

// this function is for GoRoutine
func gtpSessionRoutine(session *gtpSession) {
	myLog := log.WithField("SessionID", session.id)
	myLog.Debug("Start a GTP session goroutine")

	for msg := range session.cmdReqChan {
		myLog.Debugf("Received CMD : %v", msg)

		switch cmd := msg.(type) {
		case gsc.CreateSessionReq:
			err := procCreateSession(session, cmd, myLog)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func procCreateSession(session *gtpSession, cmd gsc.CreateSessionReq, myLog *logrus.Entry) error {
	session.status = gssCSReqSending
	seqNum := session.sgwCtrl.nextSeqNum()

	recoveryIE, err := ie.NewRecovery(0, session.sgwCtrl.recovery)
	if err != nil {
		return err
	}
	bearerContextTBCIE, err := ie.NewBearerContextToBeCreatedWithinCSReq(
		ie.BearerContextToBeCreatedWithinCSReqArg{
			Ebi:          session.ebi,
			BearerQoS:    cmd.BearerQoS,
			SgwDataFteid: session.sgwDataFTEID,
		})
	if err != nil {
		return err
	}

	csReqArg := gtpv2c.CreateSessionRequestArg{
		Imsi:             session.imsi,
		Msisdn:           session.msisdn,
		Mei:              cmd.Mei,
		Uli:              cmd.Uli,
		ServingNetwork:   session.servingNetwork,
		RatType:          session.ratType,
		Indication:       cmd.Indication,
		SgwCtrlFteid:     session.sgwCtrlFTEID,
		Apn:              session.apn,
		SelectionMode:    cmd.SelectionMode,
		PdnType:          session.pdnType,
		Paa:              session.paa,
		ApnRestriction:   cmd.ApnRestriction,
		ApnAmbr:          session.ambr,
		Pco:              cmd.Pco,
		BearerContextTBC: bearerContextTBCIE,
		Recovery:         recoveryIE,
	}

	csReq, err := gtpv2c.NewCreateSessionRequest(seqNum, csReqArg)
	if err != nil {
		return err
	}
	csReqBin := csReq.Marshal()

	// Send a CSReq packet to the PGW
	raddr := session.pgwCtrlAddr
	session.toCtrlSenderChan <- UDPpacket{raddr, csReqBin}

	var res gsc.Res
	afterChan := time.After(1 * time.Second)
	for {
		select {
		case recv := <-session.fromCtrlReceiverChan:
			myLog.Debugf("received packet from %v body: %v", recv.raddr, recv.body)

			// Ensure received packet has sent from correct PGW address
			if !recv.raddr.IP.Equal(session.pgwCtrlAddr.IP) ||
				recv.raddr.Port != session.pgwCtrlAddr.Port {
				myLog.Debugf("Received invalid GTPv2-C packet from unkown address : %v , expected : %v", recv.raddr, session.pgwCtrlAddr)
				continue
			}

			// Unmarchal received packet
			msg, _, err := gtpv2c.Unmarshal(recv.body)
			if err != nil {
				myLog.Debugf("Received invalid GTPv2-C packet")
				continue
			}
			myLog.Debugf("received GTPv2-C packet : %v", msg)

			// Ensure received packete is a Create Session Response
			csres, ok := msg.(*gtpv2c.CreateSessionResponse)
			if !ok {
				myLog.Debugf("Received packet is not a Create Session Response message.")
				continue
			}

			// Set PGW's F-TEIDs into the session
			session.pgwCtrlFTEID = csres.PgwCtrlFteid()
			session.pgwDataFTEID = csres.BearerContextCeated().PgwDataFteid()
			// Set PDN Address Allocation into the session
			session.paa = csres.Paa()

			res = gsc.Res{Code: gsc.ResOK, Msg: ""}
			goto eof
		case <-afterChan:
			myLog.Error("Timeout to wait Create Session Response")
			res = gsc.Res{Code: gsc.ResTimeout, Msg: "Timeout"}
			goto eof
		}
	}
eof:
	session.cmdResChan <- res

	return nil
}
