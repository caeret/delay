package delay

import "time"

type Option func(e *Delayer)

func Timeout(duration time.Duration) Option {
	return func(e *Delayer) {
		e.timeout = duration
	}
}

func Force() Option {
	return func(e *Delayer) {
		e.force = true
	}
}
