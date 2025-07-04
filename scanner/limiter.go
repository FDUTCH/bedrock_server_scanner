package scanner

import "time"

type Limiter interface {
	Limit()
}

type limiter struct {
	last  time.Time
	delay time.Duration
}

func (l *limiter) Limit() {
	time.Sleep(time.Until(l.last))
	l.last = l.last.Add(l.delay)
}

type NonLimiter struct{}

func (n NonLimiter) Limit() {
	// Dose nothing
}

func NewLimiter(pingsPerSecond uint64) Limiter {
	return &limiter{
		last:  time.Now(),
		delay: time.Second / time.Duration(pingsPerSecond),
	}
}
