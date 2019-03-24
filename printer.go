package log_monitoring

import (
	"context"
	"fmt"
	"io"

	tm "github.com/buger/goterm"
)

const alertBufferSize = 50

type printerImpl struct {
	lastAlerts []*Alert
	lastStats  *Stats
}

func NewPrinter() Printer {
	return &printerImpl{
		lastAlerts: make([]*Alert, 0),
		lastStats:  &Stats{},
	}
}

func (p *printerImpl) Run(ctx context.Context, alertChan <-chan *Alert, statsChan <-chan *Stats) {
	for {
		select {
		case alert := <-alertChan:
			p.processAlert(alert)
		case stats := <-statsChan:
			p.lastStats = stats
		case <-ctx.Done():
			return
		}
		p.print()
	}
}

func (p *printerImpl) processAlert(alert *Alert) {
	if len(p.lastAlerts) == alertBufferSize {
		p.lastAlerts = p.lastAlerts[1:]
	}
	p.lastAlerts = append(p.lastAlerts, alert)
}

func (p *printerImpl) print() {
	tm.Clear()

	statBox := tm.NewBox(50|tm.PCT, 50, 0)

	p.printLine(statBox, "10 second stats")
	p.printLine(statBox, "=============")

	p.printLine(statBox, "Total bytes sent: %d", p.lastStats.TotalBytes)

	p.printLine(statBox, "Hits by section")
	for _, v := range p.lastStats.TopSections {
		p.printLine(statBox, "%s: %v", v.Key, v.Value)
	}

	p.printLine(statBox, "Hits by method")
	for _, v := range p.lastStats.TopMethods {
		p.printLine(statBox, "%s: %v", v.Key, v.Value)
	}

	p.printLine(statBox, "Hits by status")
	for _, v := range p.lastStats.TopStatuses {
		p.printLine(statBox, "%s: %v", v.Key, v.Value)
	}

	alertBox := tm.NewBox(50|tm.PCT, 50, 0)
	p.printLine(alertBox, "High traffic alerts")
	p.printLine(alertBox, "=============")
	for _, v := range p.lastAlerts {
		p.printLine(alertBox, v.String())
	}

	_, _ = tm.Print(tm.MoveTo(statBox.String(), 0|tm.PCT, 0|tm.PCT), tm.MoveTo(alertBox.String(), 50|tm.PCT, 0|tm.PCT))
	tm.Flush()
}

func (p *printerImpl) printLine(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format+"\n", args...)
}
