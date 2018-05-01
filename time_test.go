package mns

import (
	"testing"
	"time"
)

func TestTimeUnixMillisecond(t *testing.T) {
	t1 := time.Date(2018, 5, 1, 12, 0, 0, 789000000, time.UTC)
	t2 := TimeUnixMillisecond(t1.Unix()*1000 + 789)
	if d := t2.Sub(t1); d != 0 {
		t.Errorf("have:%v, want:%v", d, time.Duration(0))
	}
}
