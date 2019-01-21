package delay

import (
	"time"
)

func NewDelayer(delay time.Duration, fn func(), options ...Option) *Delayer {
	e := new(Delayer)
	e.delay = delay
	e.callback = fn
	e.timeout = delay
	e.change = make(chan struct{}, 1)
	e.stop = make(chan struct{}, 1)
	for _, option := range options {
		option(e)
	}
	return e
}

type Delayer struct {
	timeout  time.Duration
	delay    time.Duration
	change   chan struct{}
	stop     chan struct{}
	force    bool
	callback func()
}

func (e *Delayer) Run() {
	last := time.Now()
	for {
		var changed bool
		func() {
			for {
				select {
				case <-e.change:
					changed = true
					if time.Now().Sub(last) >= e.timeout {
						return
					}
				case <-time.After(e.delay):
					return
				case <-e.stop:
					return
				}
			}
		}()
		select {
		case _, ok := <-e.stop:
			if !ok {
				return
			}
		default:
		}
		if (e.force || changed) && e.callback != nil {
			e.callback()
		}
		last = time.Now()
	}
}

func (e *Delayer) Fire() {
	select {
	case e.change <- struct{}{}:
	default:
	}
}

func (e *Delayer) Stop() {
	close(e.stop)
}
