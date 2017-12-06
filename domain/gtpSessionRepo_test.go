package domain

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

func TestGtpSessionsRepo_newSession(t *testing.T) {
	Init()

	theGtpSessionRepo := theSgwCtrlRepo.getCtrl(defaultSgwCtrlAddr).gtpSessionRepo
	// at the first, there should be no session.
	assert.Equal(t, 0, theGtpSessionRepo.numOfSessions())

	// add first session
	sgwCtrl := theSgwCtrlRepo.GetCtrl(defaultSgwCtrlAddr).(*SgwCtrl)
	sgwCtrlSendChan := make(chan UDPpacket)
	sgwCtrlFTEID, _ := ie.NewFteid(0, net.IPv4(127, 0, 0, 1), nil, ie.S5S8SgwGtpCIf, 0)
	sgwDataFTEID, _ := ie.NewFteid(0, net.IPv4(127, 0, 0, 1), nil, ie.S5S8SgwGtpUIf, 0)
	imsi, _ := ie.NewImsi(0, "22342345234")
	msisdn, _ := ie.NewMsisdn(0, "819012345678")
	ebi, _ := ie.NewEbi(0, 5)
	paa, _ := ie.NewPaa(0, ie.PdnTypeIPv4, net.IPv4(0, 0, 0, 0), nil)
	apn, _ := ie.NewApn(0, "apn.example.com")
	ambr, _ := ie.NewAmbr(0, 4294967, 4294967)
	ratType, _ := ie.NewRatType(0, 6)
	servingNetwork, _ := ie.NewServingNetwork(0, "440", "10")
	pdnType, _ := ie.NewPdnType(0, ie.PdnTypeIPv4)

	sid, err := theGtpSessionRepo.newSession(
		sgwCtrl,
		net.IPv4(100, 100, 100, 100),
		sgwCtrlSendChan,
		sgwCtrlFTEID, sgwDataFTEID,
		imsi, msisdn, ebi, paa, apn, ambr,
		ratType, servingNetwork, pdnType,
	)
	assert.Equal(t, SessionID(0), sid)
	assert.NoError(t, err)

	assert.Equal(t, 1, theGtpSessionRepo.numOfSessions())
	session := theGtpSessionRepo.findBySessionID(sid)
	assert.Equal(t, "22342345234", session.imsi.Value())

	// Error when same SGW-CTRL-TEID
	sid2, err := theGtpSessionRepo.newSession(
		sgwCtrl,
		net.IPv4(100, 100, 100, 100),
		sgwCtrlSendChan,
		sgwCtrlFTEID, sgwDataFTEID,
		imsi, msisdn, ebi, paa, apn, ambr,
		ratType, servingNetwork, pdnType,
	)
	assert.Error(t, err)

	assert.Equal(t, 1, theGtpSessionRepo.numOfSessions())
	session2 := theGtpSessionRepo.findBySessionID(sid2)
	assert.Equal(t, "22342345234", session2.imsi.Value())

	// get when the sid does not exist
	session = theGtpSessionRepo.findBySessionID(SessionID(2343242))
	assert.Nil(t, session)
}
