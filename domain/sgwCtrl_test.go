package domain

import (
	"net"
	"testing"
	"time"

	"github.com/craftone/gojiko/domain/apns"
	"github.com/craftone/gojiko/domain/gtpSessionCmd"
	"github.com/craftone/gojiko/gtp"
	"github.com/craftone/gojiko/gtpv2c"
	"github.com/craftone/gojiko/gtpv2c/ie"
	"github.com/stretchr/testify/assert"
)

func TestSgwCtrl_CreateSession(t *testing.T) {
	pgwIP := net.IPv4(127, 1, 1, 1)
	apn, _ := apns.NewApn("example.com", "440", "10", []net.IP{pgwIP})
	apns.TheRepo().Post(apn)

	sgwCtrl := theSgwCtrlRepo.getCtrl(defaultSgwCtrlAddr)
	resCh := make(chan error)
	imsi := "440101234567890"
	ebi := byte(5)
	go func() {
		_, err := sgwCtrl.CreateSession(
			imsi, "819012345678", "0123456789012345",
			"440", "10", "example.com", ebi,
		)
		resCh <- err
	}()

	// wait till the session is created
retry:
	session := sgwCtrl.gtpSessionRepo.findByImsiEbi(imsi, ebi)
	if session == nil {
		// fmt.Println("waiting")
		time.Sleep(50 * time.Microsecond)
		goto retry
	}
	// fmt.Printf("find the session! : %v\n", session)

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

	err := <-resCh
	assert.NoError(t, err)

	assert.True(t, session.paa.IPv4().Equal(paaIP))
	assert.Equal(t, pgwCtrlTEID, session.pgwCtrlFTEID.Teid())
	assert.Equal(t, pgwDataTEID, session.pgwDataFTEID.Teid())
}

func TestSgwCtrl_CreateSession_Timeout(t *testing.T) {
	pgwIP := net.IPv4(127, 1, 1, 1)
	apn, _ := apns.NewApn("example.com", "440", "10", []net.IP{pgwIP})
	apns.TheRepo().Post(apn)

	sgwCtrl := theSgwCtrlRepo.getCtrl(defaultSgwCtrlAddr)
	imsi := "440101234567892"
	ebi := byte(5)
	res, _ := sgwCtrl.CreateSession(
		imsi, "819012345679", "0123456789012345",
		"440", "10", "example.com", ebi,
	)

	// No Create Sessin Response and the session should be timed out.

	assert.Equal(t, *res, gtpSessionCmd.Res{Code: 500, Msg: "Timeout"})
}
