package domain

import (
	"net"
	"sync"
	"testing"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/stretchr/testify/assert"
)

func TestAbsSPgw_NextTeid(t *testing.T) {
	// single thread
	spgw := &absSPgw{}
	teid := spgw.nextTeid()
	assert.Equal(t, gtp.Teid(0), teid)
	teid = spgw.nextTeid()
	assert.Equal(t, gtp.Teid(1), teid)

	// multi thread
	spgw = &absSPgw{}
	wg := &sync.WaitGroup{}
	count := 1000
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			teid = spgw.nextTeid()
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, gtp.Teid(count), spgw.teidVal)
}

func TestAbsSPgw_NextSeqNum(t *testing.T) {
	// single thread
	spgw := &absSPgw{}
	SeqNum := spgw.nextSeqNum()
	assert.Equal(t, uint32(0), SeqNum)
	SeqNum = spgw.nextSeqNum()
	assert.Equal(t, uint32(1), SeqNum)

	// multi thread
	spgw = &absSPgw{}
	wg := &sync.WaitGroup{}
	count := 1000
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			SeqNum = spgw.nextSeqNum()
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, uint32(count), spgw.seqNum)

	// cycle
	spgw.seqNum = 0x800000 - 1
	seqNum := spgw.nextSeqNum()
	assert.Equal(t, uint32(0x800000-1), seqNum)
	seqNum = spgw.nextSeqNum()
	assert.Equal(t, uint32(0), seqNum)
}

func Test_absSPgw_findOrCreateOpSPgw(t *testing.T) {
	spgw := theSgwCtrlRepo.GetCtrl(defaultSgwCtrlAddr)

	// new OpSpgw
	opAddr := net.UDPAddr{
		IP:   net.IPv4(127, 1, 1, 1),
		Port: GtpControlPort,
	}
	opSpgw1, err := spgw.findOrCreateOpSPgw(opAddr)
	assert.NoError(t, err)

	// same OpSpgw
	opSpgw2, err := spgw.findOrCreateOpSPgw(opAddr)
	assert.NoError(t, err)
	assert.Equal(t, opSpgw1, opSpgw2)

	// another OpSpgw
	opSpgw3, err := spgw.findOrCreateOpSPgw(net.UDPAddr{
		IP:   net.IPv4(127, 1, 1, 2),
		Port: GtpControlPort,
	})
	assert.NoError(t, err)
	assert.NotEqual(t, opSpgw1, opSpgw3)
}
