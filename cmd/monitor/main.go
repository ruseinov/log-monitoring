package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/fraugster/cli"
	"github.com/papertrail/go-tail/follower"
	"github.com/ruseinov/log-monitoring"
)

const (
	defaultLogPath        = "/tmp/access.log"
	defaultRpsThreshold   = 10
	rollupInternalSeconds = 10
)

func main() {
	logPath := flag.String("logPath", defaultLogPath,
		fmt.Sprintf("specify path to the log file you want to monitor, defaults to %s", defaultLogPath))
	rpsThreshold := flag.Int64("rpsThreshold", defaultRpsThreshold,
		fmt.Sprintf("specify rps alert threshold, defaults to %d", defaultRpsThreshold))
	flag.Parse()

	reader, err := follower.New(*logPath, follower.Config{
		Whence: io.SeekEnd,
		Reopen: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	ctx := cli.Context()
	errChan := make(chan error, 1)

	parser := log_monitoring.NewGonxParser()

	alerter := log_monitoring.NewAlerter(log_monitoring.DefaultAlertIntervalSeconds)
	alertChan := alerter.Run(ctx, *rpsThreshold)

	monitor := log_monitoring.NewMonitor()
	logEntryChan, statsChan := monitor.Run(ctx, rollupInternalSeconds*time.Second)

	printer := log_monitoring.NewPrinter()

	go func() {
		printer.Run(ctx, alertChan, statsChan)
	}()

	go func() {
		for line := range reader.Lines() {
			entry, err := parser.Parse(line.String())
			if err != nil {
				errChan <- fmt.Errorf("log format is invalid: %s", err.Error())
				return
			}
			alerter.Incr()
			logEntryChan <- entry
		}

		if reader.Err() != nil {
			errChan <- reader.Err()
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("exiting log monitoring app")
	case err := <-errChan:
		log.Fatal(err)
	}

}
