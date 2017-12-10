package domain

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/craftone/gojiko/gtp"

	"github.com/craftone/gojiko/domain/gtpSessionCmd"

	"github.com/sirupsen/logrus"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

type SessionID uint32

// gtpSessionRepo manages all GTP sessions.
// A GTP session is identified by session id or
// by a tuple of SgwCtrlAddr, PgwCtrlAddr, imsi and ebi or
// by SgwCtrlFTEID.
type gtpSessionRepo struct {
	sessionsByID       map[SessionID]*gtpSession
	sessionsByCtrlTeid map[gtp.Teid]*gtpSession
	sessionsByImsiEbi  map[string]*gtpSession
	mtx4Map            sync.RWMutex

	nextSessionID SessionID
	mtx4Id        sync.RWMutex
}

func newGtpSessionRepo() *gtpSessionRepo {
	log.Info("Initialize GTP Sessions Repository")
	return &gtpSessionRepo{
		sessionsByID:       make(map[SessionID]*gtpSession),
		sessionsByCtrlTeid: make(map[gtp.Teid]*gtpSession),
		sessionsByImsiEbi:  make(map[string]*gtpSession),
	}
}

func (r *gtpSessionRepo) newSession(
	sgwCtrl *SgwCtrl,
	pgwCtrlIPv4 net.IP,
	sgwCtrlSendChan chan UDPpacket,
	sgwCtrlFTEID *ie.Fteid,
	sgwDataFTEID *ie.Fteid,
	imsi *ie.Imsi,
	msisdn *ie.Msisdn,
	ebi *ie.Ebi,
	paa *ie.Paa,
	apn *ie.Apn,
	ambr *ie.Ambr,
	ratType *ie.RatType,
	servingNetwork *ie.ServingNetwork,
	pdnType *ie.PdnType,
) (SessionID, error) {
	pgwCtrlFTEID, err := ie.NewFteid(0, pgwCtrlIPv4, nil, ie.S5S8PgwGtpCIf, 0)
	if err != nil {
		return 0, err
	}

	session := &gtpSession{
		id:     r.nextID(),
		status: gssIdle,
		mtx:    sync.RWMutex{},

		cmdReqChan:           make(chan gtpSessionCmd.Cmd),
		cmdResChan:           make(chan gtpSessionCmd.Res),
		toCtrlSenderChan:     sgwCtrlSendChan,
		fromCtrlReceiverChan: make(chan UDPpacket),
		toDataSenderChan:     make(chan UDPpacket),
		fromDataReceiverChan: make(chan UDPpacket),

		sgwCtrl:      sgwCtrl,
		sgwCtrlFTEID: sgwCtrlFTEID,
		sgwDataFTEID: sgwDataFTEID,
		pgwCtrlFTEID: pgwCtrlFTEID,

		sgwCtrlAddr: sgwCtrl.UDPAddr(),
		sgwDataAddr: sgwCtrl.getPair().UDPAddr(),
		pgwCtrlAddr: net.UDPAddr{IP: pgwCtrlIPv4, Port: GtpControlPort},
		pgwDataAddr: net.UDPAddr{IP: pgwCtrlIPv4, Port: GtpUserPort},

		imsi:           imsi,
		msisdn:         msisdn,
		ebi:            ebi,
		paa:            paa,
		apn:            apn,
		ambr:           ambr,
		ratType:        ratType,
		servingNetwork: servingNetwork,
		pdnType:        pdnType,
	}

	myLog := log.WithFields(logrus.Fields{
		"id":          session.id,
		"pgwCtrlIPv4": fmt.Sprint(pgwCtrlIPv4),
	})
	myLog.Debugf("New GTP session created")

	r.mtx4Map.Lock()
	defer r.mtx4Map.Unlock()
	if _, ok := r.sessionsByID[session.id]; ok {
		return 0, fmt.Errorf("There is already the session that have the ID : %d", session.id)
	}
	teid := session.sgwCtrlFTEID.Teid()
	if _, ok := r.sessionsByCtrlTeid[teid]; ok {
		return 0, fmt.Errorf("There is already the session that have the TEID : %d", teid)
	}
	imsiEbi := imsi.Value() + "_" + strconv.Itoa(int(ebi.Value()))
	if _, ok := r.sessionsByImsiEbi[imsiEbi]; ok {
		return 0, fmt.Errorf("There is already the session that have the IMSI : %s and EBI : %d", imsi.Value(), ebi.Value())
	}

	r.sessionsByID[session.id] = session
	r.sessionsByCtrlTeid[teid] = session
	r.sessionsByImsiEbi[imsiEbi] = session
	go gtpSessionRoutine(session)
	return session.id, nil
}

func (r *gtpSessionRepo) nextID() SessionID {
	r.mtx4Id.Lock()
	defer r.mtx4Id.Unlock()
	res := r.nextSessionID
	r.nextSessionID++
	return res
}

// findBySessionID returns nil when the id does not exist.
func (r *gtpSessionRepo) findBySessionID(id SessionID) *gtpSession {
	r.mtx4Map.RLock()
	defer r.mtx4Map.RUnlock()
	if val, ok := r.sessionsByID[id]; ok {
		return val
	}
	return nil
}

// findByTeid returns nil when the id does not exist.
func (r *gtpSessionRepo) findByTeid(teid gtp.Teid) *gtpSession {
	r.mtx4Map.RLock()
	defer r.mtx4Map.RUnlock()
	if val, ok := r.sessionsByCtrlTeid[teid]; ok {
		return val
	}
	return nil
}

func (r *gtpSessionRepo) findByImsiEbi(imsi string, ebi byte) *gtpSession {
	r.mtx4Map.RLock()
	defer r.mtx4Map.RUnlock()
	imsiEbi := imsi + "_" + strconv.Itoa(int(ebi))
	if val, ok := r.sessionsByImsiEbi[imsiEbi]; ok {
		return val
	}
	return nil
}

func (r *gtpSessionRepo) numOfSessions() int {
	return len(r.sessionsByID)
}
