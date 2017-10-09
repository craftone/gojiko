package domain

import (
	"sync"
	"testing"

	"github.com/craftone/gojiko/gtp"
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
