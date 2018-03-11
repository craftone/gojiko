package stats

import "time"

type Msg interface{}

type Uint64Msg struct {
	timestamp time.Time
	key       Key
	value     uint64
}

type TimeMsg struct {
	timestamp time.Time
	key       Key
	value     time.Time
}
