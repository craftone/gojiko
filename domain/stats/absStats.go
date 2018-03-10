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
	timeBook          map[Key]time.Time
}

func newAbsStats(ctx context.Context, chanLen int) *absStats {
	as := &absStats{
		toMsgReceiverChan: make(chan Msg, chanLen),
		int64book:         make(map[Key]int64),
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
	log := log.WithField("routine", "msgReceiver")
	log.Debug("Start MsgReceiver goroutine")
loop:
	select {
	case recMsg := <-as.toMsgReceiverChan:
		switch msg := recMsg.(type) {
		case Int64Msg:
			as.mtx.Lock()
			as.int64book[msg.key] += int64(msg.value)
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
		log.Debug("This goroutine is canceled")
	}
	log.Debug("End MsgReceiver goroutine")
}

//
// about Int64
//
func (as *absStats) SendInt64Msg(key Key, value int64) {
	msg := Int64Msg{
		timestamp: time.Now(),
		key:       key,
		value:     value,
	}
	as.ToMsgReceiverChan() <- msg
}

func (as *absStats) ReadInt64(key Key) int64 {
	as.mtx.RLock()
	defer as.mtx.RUnlock()
	if value, ok := as.int64book[key]; ok {
		return value
	}
	return 0
}

func (as *absStats) SetInt64(key Key, val int64) {
	as.mtx.Lock()
	defer as.mtx.Unlock()
	as.int64book[key] = val
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
