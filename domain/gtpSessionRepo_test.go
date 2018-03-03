package domain

import (
	"fmt"
	"net"
	"testing"

	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/stretchr/testify/assert"
)

func TestGtpSessionsRepo_newSession(t *testing.T) {
	fmt.Println(theSgwCtrlRepo)
	theGtpSessionRepo := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr).GtpSessionRepo
	// at the first, there should be no session.
	assert.Equal(t, 0, theGtpSessionRepo.NumOfSessions())

	// add first session
	sgwCtrl := theSgwCtrlRepo.GetCtrl(defaultSgwCtrlAddr).(*SgwCtrl)
	sgwCtrlSendChan := make(chan UDPpacket)

	sgwCtrlTEID := sgwCtrl.nextTeid()
	sgwDataTEID := sgwCtrl.Pair().nextTeid()
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

	assert.Equal(t, 1, theGtpSessionRepo.NumOfSessions())
	session := theGtpSessionRepo.FindBySessionID(sid)
	assert.Equal(t, "22342345234", session.Imsi())
	assert.Equal(t, "819012345678", session.Msisdn())
	assert.Equal(t, "490154203237518", session.Mei())
	assert.Equal(t, byte(5), session.Ebi())
	assert.Equal(t, net.IPv4(0, 0, 0, 0).To4(), session.Paa())
	assert.Equal(t, "apn.example.com", session.Apn())
	assert.Equal(t, "Uplink AMBR: 4294967 kbps, Downlink AMBR: 4294967 kbps", session.Ambr())
	assert.Equal(t, "EUTRAN (6)", session.RatType())
	assert.Equal(t, "MCC: 440, MNC: 10", session.ServingNetwork())
	assert.Equal(t, "IPv4", session.PdnType())

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
	assert.Equal(t, 1, theGtpSessionRepo.NumOfSessions())

	// Error when same IMSI and EBI
	sgwCtrlTEID2 := sgwCtrl.nextTeid()
	sgwDataTEID2 := sgwCtrl.Pair().nextTeid()
	sgwCtrlFTEID2, _ := ie.NewFteid(0, net.IPv4(127, 0, 0, 1), nil, ie.S5S8SgwGtpCIf, sgwCtrlTEID2)
	sgwDataFTEID2, _ := ie.NewFteid(0, net.IPv4(127, 0, 0, 1), nil, ie.S5S8SgwGtpUIf, sgwDataTEID2)
	_, err = theGtpSessionRepo.newSession(
		sgwCtrl,
		net.IPv4(100, 100, 100, 100),
		sgwCtrlSendChan,
		sgwCtrlFTEID2, sgwDataFTEID2,
		imsi1Ie, msisdnIe, meiIe, ebi1Ie, paaIe, apnIe, ambrIe,
		ratTypeIe, servingNetworkIe, pdnTypeIe,
	)
	assert.Error(t, err)
	assert.Equal(t, 1, theGtpSessionRepo.NumOfSessions())

	// No error when other SGW-CTRL-TEID and IMSI and EBI
	sid2, err := theGtpSessionRepo.newSession(
		sgwCtrl,
		net.IPv4(100, 100, 100, 100),
		sgwCtrlSendChan,
		sgwCtrlFTEID2, sgwDataFTEID2,
		imsi2Ie, msisdnIe, meiIe, ebi1Ie, paaIe, apnIe, ambrIe,
		ratTypeIe, servingNetworkIe, pdnTypeIe,
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, theGtpSessionRepo.NumOfSessions())

	// Assert to find 2nd session by SessionID
	session2 := theGtpSessionRepo.FindBySessionID(sid2)
	assert.Equal(t, imsi2, session2.imsi.Value())
	// Assert to find 2nd session by CtrlTEID
	session2ct := theGtpSessionRepo.FindByCtrlTeid(sgwCtrlTEID2)
	assert.Equal(t, imsi2, session2ct.imsi.Value())
	// Assert to find 2nd session by DataTEID
	session2dt := theGtpSessionRepo.FindByDataTeid(sgwDataTEID2)
	assert.Equal(t, imsi2, session2dt.imsi.Value())
	// Assert to find 2nd session by IMSI and EBI
	session2i := theGtpSessionRepo.FindByImsiEbi(imsi2, ebi1)
	assert.Equal(t, imsi2, session2i.imsi.Value())

	// find nil when the sid does not exist
	session = theGtpSessionRepo.FindBySessionID(SessionID(2343242))
	assert.Nil(t, session)

}
