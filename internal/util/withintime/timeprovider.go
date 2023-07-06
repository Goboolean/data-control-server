package withintime

import "time"



// It is an interface used on WithinDurationChecker
// to provide time specification logic injection on testing
type TimeProvider interface {
	Time() time.Time
}

type CurrentTime struct{}
func (t *CurrentTime) Time() time.Time { return time.Now() }

type FixedTime struct{time time.Time}
func (t *FixedTime) Time() time.Time { return t.time }
