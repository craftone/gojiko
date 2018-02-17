package domain

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

func TestGtpSessionsRepo_newSession(t *testing.T) {
	fmt.Println(theSgwCtrlRepo)
	theGtpSessionRepo := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr).gtpSessionRepo
	// at the first, there should be no session.
	assert.Equal(t, 0, theGtpSessionRepo.numOfSessions())

	// add first session
	sgwCtrl := theSgwCtrlRepo.GetCtrl(defaultSgwCtrlAddr).(*SgwCtrl)
	sgwCtrlSendChan := make(chan UDPpacket)

	sgwCtrlTEID := sgwCtrl.nextTeid()
	sgwDataTEID := sgwCtrl.getPair().nextTeid()
	sgwCtrlFTEID, _ := ie.NewFteid(0, net.IPv4(127, 0, 0, 1), nil, ie.S5S8SgwGtpCIf, sgwCtrlTEID)
	sgwDataFTEID, _ := ie.NewFteid(0, net.IPv4(127, 0, 0, 1), nil, ie.S5S8SgwGtpUIf, sgwDataTEID)

	imsi1 := "22342345234"
	imsi1Ie, _ := ie.NewImsi(0, imsi1)
	imsi2 := "22342345239"
	imsi2Ie, _ := ie.NewImsi(0, imsi2)
	msisdnIe, _ := ie.NewMsisdn(0, "819012345678")
	meiIe, _ := ie.NewMei(0, "490154203237518")
	ebi1 := byte(5)
	ebi1Ie, _ := ie.NewEbi(0, ebi1)
	ebi2 := byte(6)
	ebi2Ie, _ := ie.NewEbi(0, ebi2)
	paaIe, _ := ie.NewPaa(0, ie.PdnTypeIPv4, net.IPv4(0, 0, 0, 0), nil)
	apnIe, _ := ie.NewApn(0, "apn.example.com")
	ambrIe, _ := ie.NewAmbr(0, 4294967, 4294967)
	ratTypeIe, _ := ie.NewRatType(0, 6)
	servingNetworkIe, _ := ie.NewServingNetwork(0, "440", "10")
	pdnTypeIe, _ := ie.NewPdnType(0, ie.PdnTypeIPv4)

	sid, err := theGtpSessionRepo.newSession(
		sgwCtrl,
		net.IPv4(100, 100, 100, 100),
		sgwCtrlSendChan,
		sgwCtrlFTEID, sgwDataFTEID,
		imsi1Ie, msisdnIe, meiIe, ebi1Ie, paaIe, apnIe, ambrIe,
		ratTypeIe, servingNetworkIe, pdnTypeIe,
	)
	assert.Equal(t, SessionID(0), sid)
	assert.NoError(t, err)

	assert.Equal(t, 1, theGtpSessionRepo.numOfSessions())
	session := theGtpSessionRepo.findBySessionID(sid)
	assert.Equal(t, "22342345234", session.imsi.Value())

	// Error when same SGW-CTRL-TEID
	_, err = theGtpSessionRepo.newSession(
		sgwCtrl,
		net.IPv4(100, 100, 100, 100),
		sgwCtrlSendChan,
		sgwCtrlFTEID, sgwDataFTEID,
		imsi2Ie, msisdnIe, meiIe, ebi2Ie, paaIe, apnIe, ambrIe,
		ratTypeIe, servingNetworkIe, pdnTypeIe,
	)
	assert.Error(t, err)
	assert.Equal(t, 1, theGtpSessionRepo.numOfSessions())

	// Error when same IMSI and EBI
	sgwCtrlTEID2 := sgwCtrl.nextTeid()
	sgwCtrlFTEID2, _ := ie.NewFteid(0, net.IPv4(127, 0, 0, 1), nil, ie.S5S8SgwGtpCIf, sgwCtrlTEID2)
	_, err = theGtpSessionRepo.newSession(
		sgwCtrl,
		net.IPv4(100, 100, 100, 100),
		sgwCtrlSendChan,
		sgwCtrlFTEID2, sgwDataFTEID,
		imsi1Ie, msisdnIe, meiIe, ebi1Ie, paaIe, apnIe, ambrIe,
		ratTypeIe, servingNetworkIe, pdnTypeIe,
	)
	assert.Error(t, err)
	assert.Equal(t, 1, theGtpSessionRepo.numOfSessions())

	// No error when other SGW-CTRL-TEID and IMSI and EBI
	sid2, err := theGtpSessionRepo.newSession(
		sgwCtrl,
		net.IPv4(100, 100, 100, 100),
		sgwCtrlSendChan,
		sgwCtrlFTEID2, sgwDataFTEID,
		imsi2Ie, msisdnIe, meiIe, ebi1Ie, paaIe, apnIe, ambrIe,
		ratTypeIe, servingNetworkIe, pdnTypeIe,
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, theGtpSessionRepo.numOfSessions())

	// Assert to find 2nd session by SessionID
	session2 := theGtpSessionRepo.findBySessionID(sid2)
	assert.Equal(t, imsi2, session2.imsi.Value())
	// Assert to find 2nd session by TEID
	session2t := theGtpSessionRepo.findByTeid(sgwCtrlTEID2)
	assert.Equal(t, imsi2, session2t.imsi.Value())
	// Assert to find 2nd session by IMSI and EBI
	session2i := theGtpSessionRepo.findByImsiEbi(imsi2, ebi1)
	assert.Equal(t, imsi2, session2i.imsi.Value())

	// find nil when the sid does not exist
	session = theGtpSessionRepo.findBySessionID(SessionID(2343242))
	assert.Nil(t, session)
}
