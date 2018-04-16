package domain

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/domain/gtpv1u"
	"github.com/stretchr/testify/assert"
)

func TestSgwData_echoReceiver(t *testing.T) {
	// send ECHO Request to SGW Ctrl
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(defaultSgwCtrlAddr)
	sgwData := sgwCtrl.Pair()
	sgwDataAddr := sgwData.UDPAddr()
	conn, err := net.DialUDP("udp", nil, &sgwDataAddr)
	assert.NoError(t, err)
	seqNum := uint32(0x1234)
	echoReqBin := []byte{
		byte(1<<5 + 1<<4 + 1<<1), // Version:1, PT:1, Sequence:1
		byte(gtpv1u.EchoRequestNum),
		0, 4, // Length
		0, 0, 0, 0, // TEID is always 0
		byte(seqNum >> 8), byte(seqNum), // Sequence Number
		0, // N-PDU Number
		0, // Next Extention Header Type
	}
	conn.Write(echoReqBin)
	// receive ECHO Response
	recvBuf := make([]byte, 1024)
	n, err := conn.Read(recvBuf)
	assert.NoError(t, err)
	expectedEchoResBin := []byte{
		byte(1<<5 + 1<<4 + 1<<1), // Version:1, PT:1, Sequence:1
		byte(gtpv1u.EchoResponseNum),
		0, 6, // Length
		0, 0, 0, 0, // TEID is always 0
		byte(seqNum >> 8), byte(seqNum), // Sequence Number
		0,                  // N-PDU Number
		0,                  // Next Extention Header Type
		14,                 // RecoveryNum
		sgwData.Recovery(), // Recovery Value
	}
	// assertion
	assert.Equal(t, expectedEchoResBin, recvBuf[:n])
}
