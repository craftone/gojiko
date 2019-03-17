package domain

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/sirupsen/logrus"
)

type SessionID uint32

// GtpSessionRepo manages all GTP sessions.
// A GTP session is identified by session id or
// by a tuple of SgwCtrlAddr, PgwCtrlAddr, imsi and ebi or
// by SgwCtrlFTEID.
type GtpSessionRepo struct {
	sessionsByID       map[SessionID]*GtpSession
	sessionsByCtrlTeid map[gtp.Teid]*GtpSession
	sessionsByDataTeid map[gtp.Teid]*GtpSession
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
		sessionsByDataTeid: make(map[gtp.Teid]*GtpSession),
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
	tai *ie.Tai,
	ecgi *ie.Ecgi,
	servingNetwork *ie.ServingNetwork,
	pdnType *ie.PdnType,
) (SessionID, error) {
	pgwCtrlFTEID, err := ie.NewFteid(0, pgwCtrlIPv4, nil, ie.S5S8PgwGtpCIf, 0)
	if err != nil {
		return 0, err
	}

	session := &GtpSession{
		id:         r.nextID(),
		status:     GssNewed,
		mtx4status: sync.RWMutex{},

		receiveCSresChan: make(chan *gtpv2c.CreateSessionResponse),
		receiveMBresChan: make(chan *gtpv2c.ModifyBearerResponse),
		receiveDSresChan: make(chan *gtpv2c.DeleteSessionResponse),

		toSgwCtrlSenderChan:     sgwCtrlSendChan,
		fromSgwCtrlReceiverChan: make(chan UDPpacket, 10),
		toSgwDataSenderChan:     make(chan UDPpacket, 100),
		fromSgwDataReceiverChan: make(chan UDPpacket, 100),

		sgwCtrl:      sgwCtrl,
		sgwCtrlFTEID: sgwCtrlFTEID,
		sgwDataFTEID: sgwDataFTEID,
		pgwCtrlFTEID: pgwCtrlFTEID,

		sgwCtrlAddr: sgwCtrl.UDPAddr(),
		sgwDataAddr: sgwCtrl.Pair().UDPAddr(),
		pgwCtrlAddr: net.UDPAddr{IP: pgwCtrlIPv4, Port: GtpControlPort},
		// pgwDataAddr will be updated after receive CreateSessionResponse
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
		tai:            tai,
		ecgi:           ecgi,
	}

	log := log.WithFields(logrus.Fields{
		"id":          session.id,
		"pgwCtrlIPv4": pgwCtrlIPv4.String(),
		"imsi":        session.imsi.Value(),
		"ebi":         session.ebi.Value(),
		"msisdn":      session.msisdn.Value(),
	})
	log.Info("New GTP session created")

	r.mtx4Map.Lock()
	defer r.mtx4Map.Unlock()
	if _, ok := r.sessionsByID[session.id]; ok {
		return 0, NewDuplicateSessionError("There is already the session that have the ID : %d", session.id)
	}
	ctrlTeid := session.sgwCtrlFTEID.Teid()
	if _, ok := r.sessionsByCtrlTeid[ctrlTeid]; ok {
		return 0, NewDuplicateSessionError("There is already the session that have the SGW-CTRL-TEID : %d", ctrlTeid)
	}
	dataTeid := session.sgwDataFTEID.Teid()
	if _, ok := r.sessionsByDataTeid[dataTeid]; ok {
		return 0, NewDuplicateSessionError("There is already the session that have the SGW-DATA-TEID : %d", dataTeid)
	}
	imsiEbi := imsi.Value() + "_" + strconv.Itoa(int(ebi.Value()))
	if _, ok := r.sessionsByImsiEbi[imsiEbi]; ok {
		return 0, NewDuplicateSessionError("There is already the session that have the IMSI : %s and EBI : %d", imsi.Value(), ebi.Value())
	}

	r.sessionsByID[session.id] = session
	r.sessionsByCtrlTeid[ctrlTeid] = session
	r.sessionsByDataTeid[dataTeid] = session
	r.sessionsByImsiEbi[imsiEbi] = session
	go session.receiveCtrlPacketRoutine()
	go session.receiveDataPacketRoutine()
	return session.id, nil
}

func (r *GtpSessionRepo) deleteSession(session *GtpSession) error {
	sessionID := session.ID()
	log := log.WithField("SessionID", sessionID)
	log.Info("Delete a session record")
	for session.Status() == GssCSReqSend || session.Status() == GssDSReqSend {
		log.Debug("Wait for ready to delete")
		time.Sleep(time.Millisecond)
	}
	r.mtx4Map.Lock()
	defer r.mtx4Map.Unlock()
	session, ok := r.sessionsByID[sessionID]
	if !ok {
		return fmt.Errorf("There is no session with that id : %0X", sessionID)
	}
	delete(r.sessionsByID, sessionID)

	ctrlTeid := session.sgwCtrlFTEID.Teid()
	if _, ok := r.sessionsByCtrlTeid[ctrlTeid]; ok {
		delete(r.sessionsByCtrlTeid, ctrlTeid)
	} else {
		log.Errorf("There is no session with that SGW Ctrl F-TEID : %0X", ctrlTeid)
	}

	dataTeid := session.sgwDataFTEID.Teid()
	if _, ok := r.sessionsByDataTeid[dataTeid]; ok {
		delete(r.sessionsByDataTeid, dataTeid)
	} else {
		log.Errorf("There is no session with that SGW Data F-TEID : %0X", dataTeid)
	}

	imsiEbi := session.imsi.Value() + "_" + strconv.Itoa(int(session.ebi.Value()))
	if _, ok := r.sessionsByImsiEbi[imsiEbi]; ok {
		delete(r.sessionsByImsiEbi, imsiEbi)
	} else {
		log.Errorf("There is no session with that IMSI and EBI : %s", imsiEbi)
	}

	close(session.receiveCSresChan)        // receiveCSresChan is used by gtpSession only
	close(session.receiveDSresChan)        // receiveDSresChan is used by gtpSession only
	close(session.fromSgwCtrlReceiverChan) // the sender should care the channel is active
	close(session.toSgwDataSenderChan)     // tell the data sender to finish
	close(session.fromSgwDataReceiverChan) // the sender should care the channel is active

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

// FindByCtrlTeid returns nil when the id does not exist.
func (r *GtpSessionRepo) FindByCtrlTeid(teid gtp.Teid) *GtpSession {
	r.mtx4Map.RLock()
	defer r.mtx4Map.RUnlock()
	if val, ok := r.sessionsByCtrlTeid[teid]; ok {
		return val
	}
	return nil
}

// FindByDataTeid returns nil when the id does not exist.
func (r *GtpSessionRepo) FindByDataTeid(teid gtp.Teid) *GtpSession {
	r.mtx4Map.RLock()
	defer r.mtx4Map.RUnlock()
	if val, ok := r.sessionsByDataTeid[teid]; ok {
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
