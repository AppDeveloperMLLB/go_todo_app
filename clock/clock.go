package clock

import (
	"time"
)

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct {
	NowTime time.Time
}

func NewFixedClocker() Clocker {
	return &FixedClocker{
		NowTime: time.Date(2024, 9, 1, 12, 34, 56, 0, time.UTC),
	}
}

func CreateFixedClocker(t time.Time) Clocker {
	return &FixedClocker{
		NowTime: t,
	}
}

func (fc FixedClocker) Now() time.Time {
	return fc.NowTime
}
