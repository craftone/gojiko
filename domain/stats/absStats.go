package stats

import (
	"context"
	"sync"
	"time"
)

type absStats struct {
	mtx               sync.RWMutex
	toMsgReceiverChan chan Msg
	int64book         map[Key]int64
}

func newAbsStats(ctx context.Context, chanLen int) *absStats {
	as := &absStats{
		toMsgReceiverChan: make(chan Msg, chanLen),
		int64book:         make(map[Key]int64),
	}
	go as.msgReceiver(ctx)
	return as
}

func (as *absStats) ToMsgReceiverChan() chan Msg {
	return as.toMsgReceiverChan
}

// msgReceiver is for goroutine
func (as *absStats) msgReceiver(ctx context.Context) {
	log := log.WithField("routine", "msgReceiver")
	log.Debug("start MsgReceiver goroutine")
loop:
	select {
	case recMsg := <-as.toMsgReceiverChan:
		switch msg := recMsg.(type) {
		case Int64Msg:
			as.mtx.Lock()
			as.int64book[msg.key] += int64(msg.value)
			as.mtx.Unlock()
		default:
			log.Debugf("Received invalid message : %#v", recMsg)
		}
		goto loop
	case <-ctx.Done():
		log.Debug("This goroutine is canceled")
	}
	log.Debug("end MsgReceiver goroutine")
}

func (as *absStats) ReadInt64(key Key) (int64, bool) {
	as.mtx.RLock()
	defer as.mtx.RUnlock()
	value, ok := as.int64book[key]
	return value, ok
}

func (as *absStats) SetInt64(key Key, val int64) {
	as.mtx.Lock()
	defer as.mtx.Unlock()
	as.int64book[key] = val
}

func (as *absStats) SendInt64Msg(key Key, value int64) {
	msg := Int64Msg{
		timestamp: time.Now(),
		key:       key,
		value:     value,
	}
	as.ToMsgReceiverChan() <- msg
}
