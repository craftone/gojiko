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
	pgwCtrlFTEID *ie.Fteid
	pgwDataFTEID *ie.Fteid

	pgwCtrlAddr net.UDPAddr
	pgwDataAddr net.UDPAddr
	sgwCtrlAddr net.UDPAddr
	sgwDataAddr net.UDPAddr

	imsi           *ie.Imsi
	msisdn         *ie.Msisdn
	ebi            *ie.Ebi
	paa            *ie.Paa
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

	select {
	case recv := <-session.fromCtrlReceiverChan:
		myLog.Debugf("received packet from %v body: %v", recv.raddr, recv.body)
		session.cmdResChan <- gsc.Res{Code: gsc.ResOK, Msg: "OK?"}
	case <-time.After(1 * time.Second):
		myLog.Error("Timeout waiting Create Session Response")
		session.cmdResChan <- gsc.Res{Code: gsc.ResTimeout, Msg: "Timeout"}
	}

	return nil
}
