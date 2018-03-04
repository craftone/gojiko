package stats

import "time"

type Msg interface{}

type Int64Msg struct {
	timestamp time.Time
	key       Key
	value     int64
}
