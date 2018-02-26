package domain

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/craftone/gojiko/gtpv2c"

	"github.com/craftone/gojiko/gtp"

	"github.com/sirupsen/logrus"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

type SessionID uint32

// GtpSessionRepo manages all GTP sessions.
// A GTP session is identified by session id or
// by a tuple of SgwCtrlAddr, PgwCtrlAddr, imsi and ebi or
// by SgwCtrlFTEID.
type GtpSessionRepo struct {
	sessionsByID       map[SessionID]*GtpSession
	sessionsByCtrlTeid map[gtp.Teid]*GtpSession
	sessionsByImsiEbi  map[string]*GtpSession
	mtx4Map            sync.RWMutex

	nextSessionID SessionID
	mtx4Id        sync.RWMutex
}

func newGtpSessionRepo() *GtpSessionRepo {
	log.Info("Initialize GTP Sessions Repository")
	return &GtpSessionRepo{
		sessionsByID:       make(map[SessionID]*GtpSession),
		sessionsByCtrlTeid: make(map[gtp.Teid]*GtpSession),
		sessionsByImsiEbi:  make(map[string]*GtpSession),
	}
}

func (r *GtpSessionRepo) newSession(
	sgwCtrl *SgwCtrl,
	pgwCtrlIPv4 net.IP,
	sgwCtrlSendChan chan UDPpacket,
	sgwCtrlFTEID *ie.Fteid,
	sgwDataFTEID *ie.Fteid,
	imsi *ie.Imsi,
	msisdn *ie.Msisdn,
	mei *ie.Mei,
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

	session := &GtpSession{
		id:         r.nextID(),
		status:     GssIdle,
		mtx4status: sync.RWMutex{},

		cmdReqChan:       make(chan gtpSessionCmd, 10),
		cmdResChan:       make(chan GscRes, 10),
		receiveCSresChan: make(chan *gtpv2c.CreateSessionResponse, 10),

		toCtrlSenderChan:     sgwCtrlSendChan,
		fromCtrlReceiverChan: make(chan UDPpacket, 10),
		toDataSenderChan:     make(chan UDPpacket, 100),
		fromDataReceiverChan: make(chan UDPpacket, 100),

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
		mei:            mei,
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
		"imsi":        session.imsi,
		"ebi":         session.ebi,
		"msisdn":      session.msisdn,
	})
	myLog.Debugf("New GTP session created")

	r.mtx4Map.Lock()
	defer r.mtx4Map.Unlock()
	if _, ok := r.sessionsByID[session.id]; ok {
		return 0, fmt.Errorf("There is already the session that have the ID : %d", session.id)
	}
	teid := session.sgwCtrlFTEID.Teid()
	if _, ok := r.sessionsByCtrlTeid[teid]; ok {
		return 0, fmt.Errorf("There is already the session that have the SGW-CTRL-TEID : %d", teid)
	}
	imsiEbi := imsi.Value() + "_" + strconv.Itoa(int(ebi.Value()))
	if _, ok := r.sessionsByImsiEbi[imsiEbi]; ok {
		return 0, fmt.Errorf("There is already the session that have the IMSI : %s and EBI : %d", imsi.Value(), ebi.Value())
	}

	r.sessionsByID[session.id] = session
	r.sessionsByCtrlTeid[teid] = session
	r.sessionsByImsiEbi[imsiEbi] = session
	go session.gtpSessionRoutine()
	go session.receivePacketRoutine()
	return session.id, nil
}

func (r *GtpSessionRepo) deleteSession(sessionID SessionID) error {
	r.mtx4Map.Lock()
	defer r.mtx4Map.Unlock()
	session, ok := r.sessionsByID[sessionID]
	if !ok {
		return fmt.Errorf("There is no session with that id : %0X", sessionID)
	}
	delete(r.sessionsByID, sessionID)

	teid := session.sgwCtrlFTEID.Teid()
	if _, ok := r.sessionsByCtrlTeid[teid]; ok {
		delete(r.sessionsByCtrlTeid, teid)
	} else {
		log.Debugf("There is no session with that SGW Ctrl F-TEID : %0X", teid)
	}

	imsiEbi := session.imsi.Value() + "_" + strconv.Itoa(int(session.ebi.Value()))
	if _, ok := r.sessionsByImsiEbi[imsiEbi]; ok {
		delete(r.sessionsByImsiEbi, imsiEbi)
	} else {
		log.Debugf("There is no session with that IMSI and EBI : %s", imsiEbi)
	}

	close(session.cmdReqChan)           // tell the gtpSession to finish
	close(session.cmdResChan)           // cmdResChan is used by gtpSession only
	close(session.receiveCSresChan)     // receiveCSresChan is used by gtpSession only
	close(session.fromCtrlReceiverChan) // the sender should care the channel is active
	close(session.toDataSenderChan)     // tell the data sender to finish
	close(session.fromDataReceiverChan) // the sender should care the channel is active

	return nil
}

func (r *GtpSessionRepo) nextID() SessionID {
	r.mtx4Id.Lock()
	defer r.mtx4Id.Unlock()
	res := r.nextSessionID
	r.nextSessionID++
	return res
}

// FindBySessionID returns nil when the id does not exist.
func (r *GtpSessionRepo) FindBySessionID(id SessionID) *GtpSession {
	r.mtx4Map.RLock()
	defer r.mtx4Map.RUnlock()
	if val, ok := r.sessionsByID[id]; ok {
		return val
	}
	return nil
}

// FindByTeid returns nil when the id does not exist.
func (r *GtpSessionRepo) FindByTeid(teid gtp.Teid) *GtpSession {
	r.mtx4Map.RLock()
	defer r.mtx4Map.RUnlock()
	if val, ok := r.sessionsByCtrlTeid[teid]; ok {
		return val
	}
	return nil
}

func (r *GtpSessionRepo) FindByImsiEbi(imsi string, ebi byte) *GtpSession {
	r.mtx4Map.RLock()
	defer r.mtx4Map.RUnlock()
	imsiEbi := imsi + "_" + strconv.Itoa(int(ebi))
	if val, ok := r.sessionsByImsiEbi[imsiEbi]; ok {
		return val
	}
	return nil
}

func (r *GtpSessionRepo) NumOfSessions() int {
	return len(r.sessionsByID)
}

func (r *GtpSessionRepo) GetSessions(maxNum int) []*GtpSession {
	res := make([]*GtpSession, r.NumOfSessions())
	n := 0
	for _, v := range r.sessionsByID {
		res = append(res, v)
		n++
		if n == maxNum {
			break
		}
	}
	return res
}
