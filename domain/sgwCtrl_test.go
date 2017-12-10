package domain

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/craftone/gojiko/domain/apns"
	"github.com/stretchr/testify/assert"
)

func TestSgwCtrl_CreateSession(t *testing.T) {
	Init()

	apn, _ := apns.NewApn("example.com", "440", "10", []net.IP{net.IPv4(127, 1, 1, 1)})
	apns.TheRepo().Post(apn)

	sgwCtrl := theSgwCtrlRepo.getCtrl(defaultSgwCtrlAddr)
	resCh := make(chan error)
	imsi := "440101234567890"
	ebi := byte(5)
	go func() {
		resCh <- sgwCtrl.CreateSession(
			imsi, "819012345678", "0123456789012345",
			"440", "10", "example.com", ebi,
		)
	}()

	// wait till the session is created
retry:
	session := sgwCtrl.gtpSessionRepo.findByImsiEbi(imsi, ebi)
	if session == nil {
		// fmt.Println("waiting")
		time.Sleep(50 * time.Microsecond)
		goto retry
	}
	fmt.Printf("find! : %v\n", session)

	session.fromCtrlReceiverChan <- UDPpacket{defaultSgwCtrlAddr, []byte{0, 0, 0}}

	err := <-resCh
	assert.NoError(t, err)
}
