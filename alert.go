package log_monitoring

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	alertHighTraffic            = 0
	alertRecovered              = 1
	DefaultAlertIntervalSeconds = 120
)

func NewAlerter(alertIntervalSeconds int64) Alerter {
	if alertIntervalSeconds <= 0 {
		panic("alert interval should be positive")
	}

	return &alerterImpl{
		requestNumLock:       &sync.Mutex{},
		alertIntervalSeconds: alertIntervalSeconds,
	}
}

type alerterImpl struct {
	exceededBandwidth    bool
	alertIntervalSeconds int64
	requestNum           int64
	requestNumLock       sync.Locker
}

func (a *alerterImpl) Incr() {
	a.requestNumLock.Lock()
	defer a.requestNumLock.Unlock()
	a.requestNum += 1
}

func (a *alerterImpl) sendAlert(ch chan *Alert, rpsThreshold int64) {
	a.requestNumLock.Lock()
	defer a.requestNumLock.Unlock()

	rps := a.requestNum / a.alertIntervalSeconds
	a.requestNum = 0
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

func (a *alerterImpl) Run(ctx context.Context, rpsThreshold int64) <-chan *Alert {
	ch := make(chan *Alert, 0)

	go func() {
		tick := time.NewTicker(time.Duration(a.alertIntervalSeconds) * time.Second)
		defer tick.Stop()

		for {
			select {
			case <-tick.C:
				a.sendAlert(ch, rpsThreshold)
			case <-ctx.Done():
				log.Println("shutting down alerter")
				return
			}
		}
	}()

	return ch
}

type Alert struct {
	Kind      int
	Rps       int64
	Timestamp time.Time
}

func (a *Alert) String() string {
	if a.Kind == alertHighTraffic {
		return fmt.Sprintf("[%v] High traffic detected: %d rps", a.Timestamp, a.Rps)
	}

	return fmt.Sprintf("[%v] Traffic back to normal: %d rps", a.Timestamp, a.Rps)
}

func newAlert(rps int64, kind int) *Alert {
	return &Alert{
		Kind:      kind,
		Rps:       rps,
		Timestamp: time.Now(),
	}
}
