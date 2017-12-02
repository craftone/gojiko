package domain

import (
	"net"
	"sync"

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

	cmdChan      chan gtpSessionCmd
	ctrlSendChan chan UDPpacket
	ctrlRecvChan chan UDPpacket
	dataSendChan chan UDPpacket
	dataRecvChan chan UDPpacket

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

	for msg := range session.cmdChan {
		myLog.Debugf("Received CMD : %v", msg)

		switch cmd := msg.(type) {
		case gscCreateSession:
			err := procCreateSession(session, cmd)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func procCreateSession(session *gtpSession, cmd gscCreateSession) error {
	session.status = gssCSReqSending
	seqNum := session.sgwCtrl.nextSeqNum()

	recoveryIE, err := ie.NewRecovery(0, session.sgwCtrl.recovery)
	if err != nil {
		return err
	}
	bearerContextTBCIE, err := ie.NewBearerContextToBeCreatedWithinCSReq(
		ie.BearerContextToBeCreatedWithinCSReqArg{
			Ebi:          session.ebi,
			BearerQoS:    cmd.bearerQoS,
			SgwDataFteid: session.sgwDataFTEID,
		})
	if err != nil {
		return err
	}

	csReqArg := gtpv2c.CreateSessionRequestArg{
		Imsi:             session.imsi,
		Msisdn:           session.msisdn,
		Mei:              cmd.mei,
		Uli:              cmd.uli,
		ServingNetwork:   session.servingNetwork,
		RatType:          session.ratType,
		Indication:       cmd.indication,
		SgwCtrlFteid:     session.sgwCtrlFTEID,
		Apn:              session.apn,
		SelectionMode:    cmd.selectionMode,
		PdnType:          session.pdnType,
		Paa:              session.paa,
		ApnRestriction:   cmd.apnRestriction,
		ApnAmbr:          session.ambr,
		Pco:              cmd.pco,
		BearerContextTBC: bearerContextTBCIE,
		Recovery:         recoveryIE,
	}

	csReq, err := gtpv2c.NewCreateSessionRequest(seqNum, csReqArg)
	if err != nil {
		return err
	}
	csReqBin := csReq.Marshal()

	raddr := session.pgwCtrlAddr
	session.ctrlSendChan <- UDPpacket{raddr, csReqBin}
	return nil
}
