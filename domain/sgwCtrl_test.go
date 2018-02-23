package domain

import (
	"net"
	"testing"
	"time"

	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/gtp"
	"github.com/craftone/gojiko/gtpv2c"
	"github.com/craftone/gojiko/gtpv2c/ie"
	"github.com/stretchr/testify/assert"
)

type csResStr struct {
	res *GscRes
	err error
}

func TestSgwCtrl_CreateSession_OK(t *testing.T) {
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	resCh := make(chan csResStr)
	imsi := "440101234567890"
	ebi := byte(5)
	go func() {
		res, err := sgwCtrl.CreateSession(
			imsi, "819012345678", "0123456789012345",
			"440", "10", "example.com", ebi,
		)
		resCh <- csResStr{res, err}
	}()

	// wait till the session is created
	session := ensureTheSession(sgwCtrl, imsi, ebi)

	pgwAddr := net.UDPAddr{IP: pgwIP, Port: GtpControlPort}

	// send invalid binary
	session.fromCtrlReceiverChan <- UDPpacket{pgwAddr, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}

	// make pseudo response binary that cause is CauseRequestAccepted
	paaIP := net.IPv4(9, 10, 11, 12)
	pgwCtrlTEID := gtp.Teid(0x01234567)
	pgwDataTEID := gtp.Teid(0x76543210)
	csResArg, _ := gtpv2c.MakeCSResArg(
		session.sgwCtrlFTEID.Teid(), // SgwCtrlTEID
		ie.CauseRequestAccepted,     // Cause
		pgwIP, pgwCtrlTEID, // PGW Ctrl FTEID
		pgwIP, pgwDataTEID, // PGW Data FTEID
		paaIP,                // PDN Allocated IP address
		net.IPv4(8, 8, 8, 8), // PriDNS
		net.IPv4(8, 8, 4, 4), // SecDNS
		5)                    // EBI
	csRes, _ := gtpv2c.NewCreateSessionResponse(0x1234, csResArg)
	csResBin := csRes.Marshal()

	// send from invalid address
	session.fromCtrlReceiverChan <- UDPpacket{defaultSgwCtrlAddr, csResBin}

	// send from invalid port
	invalidPort := net.UDPAddr{IP: pgwAddr.IP, Port: 1}
	session.fromCtrlReceiverChan <- UDPpacket{invalidPort, csResBin}

	// send valid packet
	session.fromCtrlReceiverChan <- UDPpacket{pgwAddr, csResBin}

	csres := <-resCh
	assert.NoError(t, csres.err)
	assert.Equal(t, GscResOK, csres.res.Code)

	assert.True(t, session.paa.IPv4().Equal(paaIP))
	assert.Equal(t, pgwCtrlTEID, session.pgwCtrlFTEID.Teid())
	assert.Equal(t, pgwDataTEID, session.pgwDataFTEID.Teid())
}

func ensureTheSession(sgwCtrl *SgwCtrl, imsi string, ebi byte) *GtpSession {
retry:
	session := sgwCtrl.gtpSessionRepo.findByImsiEbi(imsi, ebi)
	if session == nil {
		// fmt.Println("waiting")
		time.Sleep(50 * time.Microsecond)
		goto retry
	}
	// fmt.Printf("find the session! : %v\n", session)
	return session
}

func TestSgwCtrl_CreateSession_RetryableNG(t *testing.T) {
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	resCh := make(chan csResStr)
	imsi := "440101234567891"
	ebi := byte(5)

	go func() {
		res, err := sgwCtrl.CreateSession(
			imsi, "819012345671", "0123456789012345",
			"440", "10", "example.com", ebi,
		)
		resCh <- csResStr{res, err}
	}()

	// wait till the session is created
	session := ensureTheSession(sgwCtrl, imsi, ebi)

	pgwAddr := net.UDPAddr{IP: pgwIP, Port: GtpControlPort}

	// make pseudo response binary that cause is CauseNoResourcesAvailable
	paaIP := net.IPv4(9, 10, 11, 12)
	pgwCtrlTEID := gtp.Teid(0x01234567)
	pgwDataTEID := gtp.Teid(0x76543210)
	csResArg, _ := gtpv2c.MakeCSResArg(
		session.sgwCtrlFTEID.Teid(),  // SgwCtrlTEID
		ie.CauseNoResourcesAvailable, // Cause
		pgwIP, pgwCtrlTEID, // PGW Ctrl FTEID
		pgwIP, pgwDataTEID, // PGW Data FTEID
		paaIP,                // PDN Allocated IP address
		net.IPv4(8, 8, 8, 8), // PriDNS
		net.IPv4(8, 8, 4, 4), // SecDNS
		5)                    // EBI
	csRes, _ := gtpv2c.NewCreateSessionResponse(0x1234, csResArg)
	csResBin := csRes.Marshal()

	// send Retryable NG packet
	session.fromCtrlReceiverChan <- UDPpacket{pgwAddr, csResBin}

	csres := <-resCh

	session = sgwCtrl.gtpSessionRepo.findByImsiEbi(imsi, ebi)
	assert.Nil(t, session)
	assert.NoError(t, csres.err)
	assert.Equal(t, GscResRetryableNG, csres.res.Code)
}

func TestSgwCtrl_CreateSession_NG(t *testing.T) {
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	resCh := make(chan csResStr)
	imsi := "440101234567892"
	ebi := byte(5)

	go func() {
		res, err := sgwCtrl.CreateSession(
			imsi, "819012345671", "0123456789012345",
			"440", "10", "example.com", ebi,
		)
		resCh <- csResStr{res, err}
	}()

	// wait till the session is created
	session := ensureTheSession(sgwCtrl, imsi, ebi)

	pgwAddr := net.UDPAddr{IP: pgwIP, Port: GtpControlPort}

	// make pseudo response binary that cause is CauseNoResourcesAvailable
	paaIP := net.IPv4(9, 10, 11, 12)
	pgwCtrlTEID := gtp.Teid(0x01234567)
	pgwDataTEID := gtp.Teid(0x76543210)
	csResArg, _ := gtpv2c.MakeCSResArg(
		session.sgwCtrlFTEID.Teid(), // SgwCtrlTEID
		ie.CauseIMSINotKnown,        // Cause
		pgwIP, pgwCtrlTEID,          // PGW Ctrl FTEID
		pgwIP, pgwDataTEID, // PGW Data FTEID
		paaIP,                // PDN Allocated IP address
		net.IPv4(8, 8, 8, 8), // PriDNS
		net.IPv4(8, 8, 4, 4), // SecDNS
		5)                    // EBI
	csRes, _ := gtpv2c.NewCreateSessionResponse(0x1234, csResArg)
	csResBin := csRes.Marshal()

	// send Retryable NG packet
	session.fromCtrlReceiverChan <- UDPpacket{pgwAddr, csResBin}

	csres := <-resCh

	session = sgwCtrl.gtpSessionRepo.findByImsiEbi(imsi, ebi)
	assert.Nil(t, session)
	assert.NoError(t, csres.err)
	assert.Equal(t, GscResNG, csres.res.Code)
}

func TestSgwCtrl_CreateSession_Timeout(t *testing.T) {
	// change Gtpv2cTimeout temporarily
	defaultGtpv2cTimeout := config.Gtpv2cTimeout()
	config.SetGtpv2cTimeout(1)
	defer config.SetGtpv2cTimeout(defaultGtpv2cTimeout)

	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	imsi := "440101234567892"
	ebi := byte(5)
	res, _ := sgwCtrl.CreateSession(
		imsi, "819012345679", "0123456789012345",
		"440", "10", "example.com", ebi,
	)

	// No Create Sessin Response and the session should be timed out.

	assert.Equal(t, *res, GscRes{Code: GscResTimeout, Msg: "Timeout"})
}

func TestSgwCtrl_EchoResponse(t *testing.T) {
	// preparing
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	toSender := make(chan UDPpacket, 10)
	sgwCtrl.toSender = toSender

	// receive valid echo-request
	echoReq, _ := gtpv2c.NewEchoRequest(1, 1)
	echoReqBin := echoReq.Marshal()
	raddr := net.UDPAddr{IP: net.IPv4(127, 0, 0, 2), Port: GtpControlPort}
	udpPacket := UDPpacket{raddr, echoReqBin}
	sgwCtrl.toEchoReceiver <- udpPacket

	// send echo-response
	sendPacket := <-toSender
	assert.Equal(t, raddr, sendPacket.raddr)
	echoRes, _ := gtpv2c.NewEchoResponse(sgwCtrl.seqNum-1, sgwCtrl.recovery)
	assert.Equal(t, echoRes.Marshal(), sendPacket.body)
}
