package log_monitoring

const	(
	// alertBufferSize is a number of alerts to keep in-memory and on-screen
 	alertBufferSize = 50
	alertHighTraffic = 0
	alertRecovered = 1
	alertIntervalSeconds = 160
)

type MonitorImpl struct {
	rpsThreshold int64
	logPath string
}

type Monitor interface {
	Run()
}