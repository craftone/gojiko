package domain

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
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
	session.fromSgwCtrlReceiverChan <- UDPpacket{pgwAddr, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}

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
	session.fromSgwCtrlReceiverChan <- UDPpacket{defaultSgwCtrlAddr, csResBin}

	// send from invalid port
	invalidPort := net.UDPAddr{IP: pgwAddr.IP, Port: 1}
	session.fromSgwCtrlReceiverChan <- UDPpacket{invalidPort, csResBin}

	// send valid packet
	session.fromSgwCtrlReceiverChan <- UDPpacket{pgwAddr, csResBin}

	csres := <-resCh
	assert.NoError(t, csres.err)
	assert.Equal(t, GscResOK, csres.res.Code)

	assert.True(t, session.paa.IPv4().Equal(paaIP))
	assert.Equal(t, pgwCtrlTEID, session.pgwCtrlFTEID.Teid())
	assert.Equal(t, pgwDataTEID, session.pgwDataFTEID.Teid())
}

func ensureTheSession(sgwCtrl *SgwCtrl, imsi string, ebi byte) *GtpSession {
retry:
	session := sgwCtrl.GtpSessionRepo.FindByImsiEbi(imsi, ebi)
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
	session.fromSgwCtrlReceiverChan <- UDPpacket{pgwAddr, csResBin}

	csres := <-resCh

	session = sgwCtrl.GtpSessionRepo.FindByImsiEbi(imsi, ebi)
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
	session.fromSgwCtrlReceiverChan <- UDPpacket{pgwAddr, csResBin}

	csres := <-resCh

	session = sgwCtrl.GtpSessionRepo.FindByImsiEbi(imsi, ebi)
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

	assert.Equal(t, *res, GscRes{Code: GscResTimeout, Msg: "Create Session Request timed out"})
}

func TestSgwCtrl_EchoResponse(t *testing.T) {
	// preparing
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	toSender := make(chan UDPpacket, 10)
	sgwCtrl.toSender = toSender

	// receive valid echo-request
	echoReq, _ := gtpv2c.NewEchoRequest(100, 1)
	echoReqBin := echoReq.Marshal()
	raddr := net.UDPAddr{IP: net.IPv4(127, 0, 0, 2), Port: GtpControlPort}
	udpPacket := UDPpacket{raddr, echoReqBin}
	sgwCtrl.toEchoReceiver <- udpPacket

	// send echo-response
	sendPacket := <-toSender
	assert.Equal(t, raddr, sendPacket.raddr)
	echoRes, _ := gtpv2c.NewEchoResponse(100, sgwCtrl.recovery)
	assert.Equal(t, echoRes.Marshal(), sendPacket.body)
}

func TestSgwCtrl_CreateSessionAndDeleteBearer(t *testing.T) {
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	resCh := make(chan csResStr)
	imsi := "440101234567891"
	ebi := byte(5)
	go func() {
		res, err := sgwCtrl.CreateSession(
			imsi, "819012345679", "0123456789012345",
			"440", "10", "example.com", ebi,
		)
		resCh <- csResStr{res, err}
	}()

	// wait till the session is created
	session := ensureTheSession(sgwCtrl, imsi, ebi)

	pgwAddr := net.UDPAddr{IP: pgwIP, Port: GtpControlPort}

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

	// send valid packet
	session.fromSgwCtrlReceiverChan <- UDPpacket{pgwAddr, csResBin}

	csres := <-resCh
	assert.NoError(t, csres.err)
	assert.Equal(t, GscResOK, csres.res.Code)

	assert.True(t, session.paa.IPv4().Equal(paaIP))
	assert.Equal(t, pgwCtrlTEID, session.pgwCtrlFTEID.Teid())
	assert.Equal(t, pgwDataTEID, session.pgwDataFTEID.Teid())

	//
	// Delete Bearer Test
	//
	dbReq, _ := gtpv2c.NewDeleteBearerRequest(pgwCtrlTEID, 100, ebi)
	packet := UDPpacket{raddr: pgwAddr, body: dbReq.Marshal()}
	session.fromSgwCtrlReceiverChan <- packet

	err := ensureNoSession(sgwCtrl.GtpSessionRepo, session.ID(), 10)
	assert.NoError(t, err)
}

func ensureNoSession(repo *GtpSessionRepo, id SessionID, retryCnt int) error {
retry:
	retryCnt--
	if retryCnt == 0 {
		return fmt.Errorf("The session %d exists", id)
	}
	session := repo.FindBySessionID(id)
	if session != nil {
		fmt.Println("waiting")
		time.Sleep(50 * time.Microsecond)
		goto retry
	}
	// fmt.Printf("find the session! : %v\n", session)
	return nil
}

func TestSgwCtrl_CreateSessionAndStartUdpFlow(t *testing.T) {
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	resCh := make(chan csResStr)
	imsi := "440101234567894"
	ebi := byte(5)
	go func() {
		res, err := sgwCtrl.CreateSession(
			imsi, "819012345674", "0123456789012345",
			"440", "10", "example.com", ebi,
		)
		resCh <- csResStr{res, err}
	}()

	// wait till the session is created
	session := ensureTheSession(sgwCtrl, imsi, ebi)

	pgwAddr := net.UDPAddr{IP: pgwIP, Port: GtpControlPort}

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

	// send valid packet
	session.fromSgwCtrlReceiverChan <- UDPpacket{pgwAddr, csResBin}

	csres := <-resCh
	assert.NoError(t, csres.err)
	assert.Equal(t, GscResOK, csres.res.Code)

	assert.True(t, session.paa.IPv4().Equal(paaIP))
	assert.Equal(t, pgwCtrlTEID, session.pgwCtrlFTEID.Teid())
	assert.Equal(t, pgwDataTEID, session.pgwDataFTEID.Teid())

	// NewUdpFlow
	udpFlow := UdpEchoFlowArg{
		DestAddr:       net.UDPAddr{IP: net.IPv4(100, 100, 100, 100), Port: 10000},
		SourcePort:     10001,
		SendPacketSize: 38,
		Tos:            0,
		Ttl:            255,
		TargetBps:      15000, // interval = 38*8 / 15000 = 0.02026
		NumOfSend:      5,
		RecvPacketSize: 1450,
	}

	sgwData := session.sgwCtrl.Pair().(*SgwData)
	c := make(chan UDPpacket)
	sgwData.toSender = c

	// assert no udpFlow and lastUDPFlow
	udpFlowPointer, ok := session.UDPFlow()
	assert.Nil(t, udpFlowPointer)
	assert.False(t, ok)
	lastUDPFlowPointer, ok := session.LastUDPFlow()
	assert.Nil(t, lastUDPFlowPointer)
	assert.False(t, ok)

	// New UDPFlow
	err := session.NewUdpFlow(udpFlow)
	assert.NoError(t, err)

	packet := <-c // @ 0.00 s
	packet = <-c  // @ 0.02 s
	packet = <-c  // @ 0.04 s
	packet = <-c  // @ 0.06 s
	packet = <-c  // @ 0.08 s

	// assert udpFlow exists and no lastUDPFlow
	udpFlowPointer, ok = session.UDPFlow()
	assert.NotNil(t, udpFlowPointer)
	assert.True(t, ok)
	lastUDPFlowPointer, ok = session.LastUDPFlow()
	assert.Nil(t, lastUDPFlowPointer)
	assert.False(t, ok)

	assert.Equal(t, []byte{
		0x30,     // GTP version:1, PT=1(GTP), all flags are 0
		0xFF,     // GTP_TPDU_MSG (0xFF)
		0x00, 38, // totalLen: 38
		0x76, 0x54, 0x32, 0x10, // teid

		0x45,       // version: 4, ihl: 5
		0x00,       // tos: 0,
		0x00, 0x26, // totalLen : 38
		0x00, 0x00, // id:0
		0x40, 0x00, // fragment: 0x4000
		0xff,       // ttl
		0x11,       // protocol
		0x9e, 0xe8, // checksum
		9, 10, 11, 12, //source address
		100, 100, 100, 100, //destination address

		0x27, 0x11, // source port : 10001
		0x27, 0x10, // destination port : 10000
		0x00, 0x12, // udp size : 8+10
		0x68, 0x11, // checksum
		0x05, 0xaa, // receive udp packet size : 1450
		0, 0, 0, 0, 0, 0, 0, 5, // seqNum : 5
	}, packet.body)

	// assert no udpFlow and lastUDPFlow exits
	// wait till UDPFlow will done.
	for {
		udpFlowPointer, ok = session.UDPFlow()
		if udpFlowPointer == nil {
			break
		}
		time.Sleep(time.Microsecond)
	}
	udpFlowPointer, ok = session.UDPFlow()
	assert.Nil(t, udpFlowPointer)
	assert.False(t, ok)
	lastUDPFlowPointer, ok = session.LastUDPFlow()
	assert.NotNil(t, lastUDPFlowPointer)
	assert.True(t, ok)

	//
	// Delete Bearer Test
	//
	dbReq, _ := gtpv2c.NewDeleteBearerRequest(pgwCtrlTEID, 100, ebi)
	packet = UDPpacket{raddr: pgwAddr, body: dbReq.Marshal()}
	session.fromSgwCtrlReceiverChan <- packet

	err = ensureNoSession(sgwCtrl.GtpSessionRepo, session.ID(), 10)
	assert.NoError(t, err)
}

func TestSgwCtrl_Create2SessionsAndStartUdpFlow(t *testing.T) {
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)

	resCh1 := make(chan csResStr)
	resCh2 := make(chan csResStr)
	imsi1 := "440101234567894"
	imsi2 := "440101234567895"
	msisdn1 := "819012345674"
	msisdn2 := "819012345675"
	ebi := byte(5)

	go func() {
		res1, err := sgwCtrl.CreateSession(
			imsi1, msisdn1, "0123456789012345",
			"440", "10", "example.com", ebi,
		)
		resCh1 <- csResStr{res1, err}
	}()
	go func() {
		res2, err := sgwCtrl.CreateSession(
			imsi2, msisdn2, "0123456789012345",
			"440", "10", "example.com", ebi,
		)
		resCh2 <- csResStr{res2, err}
	}()

	// wait till the session is created
	session1 := ensureTheSession(sgwCtrl, imsi1, ebi)
	session2 := ensureTheSession(sgwCtrl, imsi2, ebi)

	pgwAddr := net.UDPAddr{IP: pgwIP, Port: GtpControlPort}

	// make pseudo response binary that cause is CauseRequestAccepted
	paaIP1 := net.IPv4(11, 11, 11, 11)
	paaIP2 := net.IPv4(22, 22, 22, 22)
	pgwCtrlTEID1 := gtp.Teid(0x01234567)
	pgwDataTEID1 := gtp.Teid(0x01234567)
	pgwCtrlTEID2 := gtp.Teid(0x76543210)
	pgwDataTEID2 := gtp.Teid(0x76543210)
	csResArg1, _ := gtpv2c.MakeCSResArg(
		session1.sgwCtrlFTEID.Teid(), // SgwCtrlTEID
		ie.CauseRequestAccepted,      // Cause
		pgwIP, pgwCtrlTEID1, // PGW Ctrl FTEID
		pgwIP, pgwDataTEID1, // PGW Data FTEID
		paaIP1,               // PDN Allocated IP address
		net.IPv4(8, 8, 8, 8), // PriDNS
		net.IPv4(8, 8, 4, 4), // SecDNS
		5)                    // EBI
	csRes1, _ := gtpv2c.NewCreateSessionResponse(0x1234, csResArg1)
	csRes1Bin := csRes1.Marshal()
	csResArg2, _ := gtpv2c.MakeCSResArg(
		session2.sgwCtrlFTEID.Teid(), // SgwCtrlTEID
		ie.CauseRequestAccepted,      // Cause
		pgwIP, pgwCtrlTEID2, // PGW Ctrl FTEID
		pgwIP, pgwDataTEID2, // PGW Data FTEID
		paaIP2,               // PDN Allocated IP address
		net.IPv4(8, 8, 8, 8), // PriDNS
		net.IPv4(8, 8, 4, 4), // SecDNS
		5)                    // EBI
	csRes2, _ := gtpv2c.NewCreateSessionResponse(0x1234, csResArg2)
	csRes2Bin := csRes2.Marshal()

	// send valid packet
	session1.fromSgwCtrlReceiverChan <- UDPpacket{pgwAddr, csRes1Bin}
	session2.fromSgwCtrlReceiverChan <- UDPpacket{pgwAddr, csRes2Bin}

	csres1 := <-resCh1
	assert.NoError(t, csres1.err)
	assert.Equal(t, GscResOK, csres1.res.Code)
	csres2 := <-resCh2
	assert.NoError(t, csres2.err)
	assert.Equal(t, GscResOK, csres2.res.Code)

	assert.True(t, session1.paa.IPv4().Equal(paaIP1))
	assert.Equal(t, pgwCtrlTEID1, session1.pgwCtrlFTEID.Teid())
	assert.Equal(t, pgwDataTEID1, session1.pgwDataFTEID.Teid())
	assert.True(t, session2.paa.IPv4().Equal(paaIP2))
	assert.Equal(t, pgwCtrlTEID2, session2.pgwCtrlFTEID.Teid())
	assert.Equal(t, pgwDataTEID2, session2.pgwDataFTEID.Teid())

	// NewUdpFlow
	udpFlow1 := UdpEchoFlowArg{
		DestAddr:       net.UDPAddr{IP: net.IPv4(100, 100, 100, 100), Port: 10000},
		SourcePort:     10001,
		SendPacketSize: 38,
		Tos:            0,
		Ttl:            255,
		TargetBps:      15000, // interval = 38*8 / 15000 = 0.02026
		NumOfSend:      4,
		RecvPacketSize: 1450,
	}
	udpFlow2 := UdpEchoFlowArg{
		DestAddr:       net.UDPAddr{IP: net.IPv4(100, 100, 100, 100), Port: 10000},
		SourcePort:     10001,
		SendPacketSize: 38,
		Tos:            0,
		Ttl:            255,
		TargetBps:      15000, // interval = 38*8 / 15000 = 0.02026
		NumOfSend:      2,
		RecvPacketSize: 1450,
	}

	sgwData := session1.sgwCtrl.Pair().(*SgwData)
	c := make(chan UDPpacket)
	sgwData.toSender = c

	err := session1.NewUdpFlow(udpFlow1)
	assert.NoError(t, err)
	err = session2.NewUdpFlow(udpFlow2)
	assert.NoError(t, err)

	packet := <-c
	packet = <-c
	packet = <-c
	packet = <-c
	packet = <-c
	packet = <-c

	assert.Equal(t, []byte{
		0x30,     // GTP version:1, PT=1(GTP), all flags are 0
		0xFF,     // GTP_TPDU_MSG (0xFF)
		0x00, 38, // totalLen: 38
		0x01, 0x23, 0x45, 0x67, // teid

		0x45,       // version: 4, ihl: 5
		0x00,       // tos: 0,
		0x00, 0x26, // totalLen : 38
		0x00, 0x00, // id:0
		0x40, 0x00, // fragment: 0x4000
		0xff,       // ttl
		0x11,       // protocol
		0x9c, 0xe8, // checksum
		11, 11, 11, 11, //source address
		100, 100, 100, 100, //destination address

		0x27, 0x11, // source port : 10001
		0x27, 0x10, // destination port : 10000
		0x00, 0x12, // udp size : 8+10
		0x6a, 0x12, // checksum
		0x05, 0xaa, // receive udp packet size : 1450
		0, 0, 0, 0, 0, 0, 0, 4, // seqNum : 5
	}, packet.body)

	//
	// Delete Bearer
	//
	dbReq1, _ := gtpv2c.NewDeleteBearerRequest(pgwCtrlTEID1, 100, ebi)
	packet = UDPpacket{raddr: pgwAddr, body: dbReq1.Marshal()}
	session1.fromSgwCtrlReceiverChan <- packet
	err = ensureNoSession(sgwCtrl.GtpSessionRepo, session1.ID(), 10)
	assert.NoError(t, err)

	dbReq2, _ := gtpv2c.NewDeleteBearerRequest(pgwCtrlTEID2, 100, ebi)
	packet = UDPpacket{raddr: pgwAddr, body: dbReq2.Marshal()}
	session2.fromSgwCtrlReceiverChan <- packet
	err = ensureNoSession(sgwCtrl.GtpSessionRepo, session2.ID(), 10)
	assert.NoError(t, err)
	// assert.True(t, false) //dummy
}
