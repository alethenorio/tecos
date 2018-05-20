// A wrapper for the Go time.Ticker to allow for testing independent of time
package gitmono

import "time"

type Ticker interface {
	Stop()
	TickerChan() <-chan time.Time
}

type RealTicker struct {
	ticker *time.Ticker
}

func NewTicker(d time.Duration) *RealTicker {
	return &RealTicker{
		ticker: time.NewTicker(d),
	}
}

func (t *RealTicker) Stop() {
	t.ticker.Stop()
}

func (t *RealTicker) TickerChan() <-chan time.Time {
	return t.ticker.C
}
