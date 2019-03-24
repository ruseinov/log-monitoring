package log_monitoring_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ruseinov/log-monitoring"
)

func TestNewAlerterPanic(t *testing.T) {
	panicCnt := 0
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic")
		}
		panicCnt++
	}()
	_ = log_monitoring.NewAlerter(0)
	_ = log_monitoring.NewAlerter(-1)

	if panicCnt < 2 {
		t.Fatal("Expected 2 panics")
	}
}

func TestAlerts(t *testing.T) {
	alertInterval := int64(1)
	alerter := log_monitoring.NewAlerter(alertInterval)
	ch := alerter.Run(context.Background(), 2)

	// test no alert
	alerter.Incr()
	timer := time.NewTimer(1 * time.Second)
	select {
	case <-ch:
		t.Fatal("No alert expected")
	case <-timer.C:

	}

	// test threshold exceeded
	alerter.Incr()
	alerter.Incr()

	timer = time.NewTimer(1 * time.Second)
	select {
	case alert := <-ch:
		if !strings.Contains(alert.String(), "High") {
			t.Fatal("High traffic alert expected")
		}
	case <-timer.C:
		t.Fatal("Alert expected")
	}

	// test back to normal
	timer = time.NewTimer(1 * time.Second)
	select {
	case alert := <-ch:
		if !strings.Contains(alert.String(), "normal") {
			t.Fatal("Traffic back to normal alert expected")
		}
	case <-timer.C:
		t.Fatal("Alert expected")
	}
}
