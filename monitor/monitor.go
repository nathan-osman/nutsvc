package monitor

import (
	"fmt"

	"github.com/nathan-osman/nutclient"
	"github.com/nathan-osman/nutsvc/logger"
)

// Monitor connects to a NUT endpoint and monitors it.
type Monitor struct {
	client *nutclient.Client
}

// New creates a new Monitor instance.
func New(l *logger.Logger, addr, name string) *Monitor {
	c := nutclient.New(&nutclient.Config{
		Addr: addr,
		Name: name,
		ConnectedFn: func() {
			l.Info(
				logger.EventMonitorStatus,
				fmt.Sprintf("connected to %s", addr),
			)
		},
		DisconnectedFn: func() {
			l.Warning(
				logger.EventMonitorStatus,
				"disconnected from NUT server",
			)
		},
		PowerLostFn: func() {
			l.Error(
				logger.EventMonitorStatus,
				"line power lost; UPS on battery power",
			)
		},
		PowerRestoredFn: func() {
			l.Info(
				logger.EventMonitorStatus,
				"line power restored",
			)
		},
	})
	return &Monitor{
		client: c,
	}
}

// Close shuts down the monitor.
func (m *Monitor) Close() {
	m.client.Close()
}
