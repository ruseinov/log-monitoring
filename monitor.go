package log_monitoring

import (
	"context"
	"log"
	"math"
	"sort"
	"time"
)

const topLimit = 10

type monitorImpl struct {
	requestNumByStatus  map[string]int64
	requestNumBySection map[string]int64
	requestNumByMethod  map[string]int64
	totalBytes          int64
}

func NewMonitor() Monitor {
	impl := &monitorImpl{}
	impl.resetStats()
	return impl
}

func (m *monitorImpl) Run(ctx context.Context, rollupInterval time.Duration) (chan<- *LogEntry, <-chan *Stats) {
	entryChan := make(chan *LogEntry, 100)
	statsChan := make(chan *Stats, 10)

	go func() {
		tick := time.NewTicker(rollupInterval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("shutting down monitor")
				return
			case entry := <-entryChan:
				m.processEntry(entry)
			case <-tick.C:
				statsChan <- m.stats()
				m.resetStats()
			}
		}
	}()

	return entryChan, statsChan
}

func (m *monitorImpl) stats() *Stats {
	return &Stats{
		TopSections: m.processCounters(m.requestNumBySection),
		TopStatuses: m.processCounters(m.requestNumByStatus),
		TopMethods:  m.processCounters(m.requestNumByMethod),
		TotalBytes:  m.totalBytes,
	}
}

func (m *monitorImpl) processCounters(counters map[string]int64) []Counter {
	var sortedCounters []Counter

	for k, v := range counters {
		cntr := Counter{
			Key:   k,
			Value: v,
		}
		sortedCounters = append(sortedCounters, cntr)
	}

	sort.Slice(sortedCounters, func(prev, next int) bool {
		return sortedCounters[prev].Value > sortedCounters[next].Value
	})

	len := math.Min(topLimit, float64(len(sortedCounters)))

	return sortedCounters[:int(len)]
}

func (m *monitorImpl) resetStats() {
	m.requestNumByMethod = make(map[string]int64)
	m.requestNumBySection = make(map[string]int64)
	m.requestNumByStatus = make(map[string]int64)
	m.totalBytes = 0
}

func (m *monitorImpl) processEntry(entry *LogEntry) {
	m.totalBytes += int64(entry.Bytes)
	m.requestNumByStatus[entry.Status]++
	m.requestNumBySection[entry.Section]++
	m.requestNumByMethod[entry.Method]++
}
