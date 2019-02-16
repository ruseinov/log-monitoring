package log_monitoring

import (
	"context"
	"sync"
	"time"
)

type Alerter interface {
	Incr()
	Run(ctx context.Context, rpsThreshold int64, alertInterval time.Duration) chan <-alert
}

func newAlerter() Alerter {
	return &alerterImpl{
		requestNumLock: &sync.Mutex{},
	}
}

type alerterImpl struct {
	exceededBandwidth bool
	requestNum int64
	requestNumLock sync.Locker
}

func(a *alerterImpl) Incr() {
	a.requestNumLock.Lock()
	defer a.requestNumLock.Unlock()
	a.requestNum += 1
}

func(a *alerterImpl) sendAlert(ch chan alert, rpsThreshold int64) {
	a.requestNumLock.Lock()
	defer a.requestNumLock.Unlock()

	rps := a.requestNum / alertIntervalSeconds
	if rps > rpsThreshold {
		a.exceededBandwidth = true
		ch <- newAlert(rps, alertHighTraffic)
		return
	}

	if a.exceededBandwidth == true {
		a.exceededBandwidth = false
		ch <- newAlert(rps, alertRecovered)
	}
}

func(a *alerterImpl) Run(ctx context.Context, rpsThreshold int64, alertInterval time.Duration) chan <-alert {
	ch := make(chan alert, 0)

	go func() {
		tick := time.NewTicker(alertInterval)
		defer tick.Stop()

		for {
			select {
			case <-tick.C:
				a.sendAlert(ch, rpsThreshold)
			case <- ctx.Done():
				return
			}
		}
	}()

	return ch
}

type alert struct {
	kind int
	rps int64
	time time.Time
}

func newAlert(rps int64, kind int) alert {
	return alert{
		kind: kind,
		rps:  rps,
		time: time.Now(),
	}
}