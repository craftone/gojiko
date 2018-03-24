package stats

import (
	"context"
	"sync"
	"time"
)

type absStats struct {
	mtx               sync.RWMutex
	toMsgReceiverChan chan Msg
	uint64book        map[Key]uint64
	timeBook          map[Key]time.Time
}

func newAbsStats(ctx context.Context, chanLen int) *absStats {
	as := &absStats{
		toMsgReceiverChan: make(chan Msg, chanLen),
		uint64book:        make(map[Key]uint64),
		timeBook:          make(map[Key]time.Time),
	}
	go as.msgReceiver(ctx)
	return as
}

func (as *absStats) ToMsgReceiverChan() chan Msg {
	return as.toMsgReceiverChan
}

// msgReceiver is for goroutine
func (as *absStats) msgReceiver(ctx context.Context) {
	log := log.WithField("routine", "absStats.msgReceiver")
	log.Debug("Start MsgReceiver goroutine")
loop:
	select {
	case recMsg := <-as.toMsgReceiverChan:
		switch msg := recMsg.(type) {
		case Uint64Msg:
			as.mtx.Lock()
			as.uint64book[msg.key] += msg.value
			as.mtx.Unlock()
		case TimeMsg:
			as.mtx.Lock()
			as.timeBook[msg.key] = msg.value
			as.mtx.Unlock()
		default:
			log.Debugf("Received invalid message : %#v", recMsg)
		}
		goto loop
	case <-ctx.Done():
		// log.Debug("This goroutine is canceled")
	}
	log.Debug("End MsgReceiver goroutine")
}

//
// about Uint64
//
func (as *absStats) SendUint64Msg(key Key, value uint64) {
	msg := Uint64Msg{
		timestamp: time.Now(),
		key:       key,
		value:     value,
	}
	as.ToMsgReceiverChan() <- msg
}

func (as *absStats) ReadUint64(key Key) uint64 {
	as.mtx.RLock()
	defer as.mtx.RUnlock()
	if value, ok := as.uint64book[key]; ok {
		return value
	}
	return 0
}

func (as *absStats) SetUint64(key Key, val uint64) {
	as.mtx.Lock()
	defer as.mtx.Unlock()
	as.uint64book[key] = val
}

//
// about time.Time
//

func (as *absStats) SendTimeMsg(key Key, value time.Time) {
	msg := TimeMsg{
		timestamp: time.Now(),
		key:       key,
		value:     value,
	}
	as.ToMsgReceiverChan() <- msg
}

func (as *absStats) ReadTime(key Key) time.Time {
	as.mtx.RLock()
	defer as.mtx.RUnlock()
	if value, ok := as.timeBook[key]; ok {
		return value
	}
	return time.Time{}
}

func (as *absStats) SetTime(key Key, val time.Time) {
	as.mtx.Lock()
	defer as.mtx.Unlock()
	as.timeBook[key] = val
}
