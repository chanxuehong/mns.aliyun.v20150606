package mns

import "time"

func TimeUnixMillisecond(ms int64) time.Time {
	s := ms / 1000
	ns := (ms % 1000) * 1e6
	return time.Unix(s, ns)
}
