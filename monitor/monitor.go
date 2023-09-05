package monitor

import (
	"fmt"
	"time"

	"github.com/nathan-osman/nutclient"
	"github.com/nathan-osman/nutsvc/conf"
	"github.com/nathan-osman/nutsvc/logger"
)

const (

	// KeyWaitSeconds indicates how long to wait after power loss before
	// initiating the specified action. The default is two minutes.
	KeyWaitSeconds = "wait_seconds"

	// KeyAction is used to store the intended action after the wait period.
	KeyAction = "action"

	// ActionNothing doesn't do anything.
	ActionNothing = "nothing"

	// ActionShutdown causes a system shutdown to occur.
	ActionShutdown = "shutdown"

	// ActionHibernate causes the system to enter hibernation.
	ActionHibernate = "hibernate"
)

var defaultWaitSeconds = 120

// Monitor connects to a NUT endpoint and monitors it.
type Monitor struct {
	addr           string
	name           string
	client         *nutclient.Client
	logger         *logger.Logger
	conf           *conf.Conf
	startTimerChan chan any
	stopTimerChan  chan any
	closeChan      chan any
	closedChan     chan any
}

func (m *Monitor) doAction() {

	// Load the action
	v, err := m.conf.Get(KeyAction)
	if err != nil {
		v = ActionNothing
	}

	// Indicate the action that is about to take place
	m.logger.Info(
		logger.EventMonitorStatus,
		fmt.Sprintf(
			"wait interval is up; performing action: \"%s\"",
			v,
		),
	)

	// Perform the action
	var apiErr error
	switch v {
	case ActionShutdown:
		apiErr = shutdown()
	case ActionHibernate:
		apiErr = hibernate()
	}

	// Log an error if something went wrong
	if apiErr != nil {
		m.logger.Error(
			logger.EventMonitorStatus,
			fmt.Sprintf(
				"error performing action: %s",
				err.Error(),
			),
		)
	}
}

func (m *Monitor) run() {
	defer close(m.closedChan)
	var timer *time.Timer
	for {

		// Load the timer channel (if timer is active)
		var timerChan <-chan time.Time
		if timer != nil {
			timerChan = timer.C
		}

		select {
		case <-timerChan:
			m.doAction()
		case <-m.startTimerChan:

			// Retrieve the wait interval
			v, err := m.conf.GetInt(KeyWaitSeconds)
			if err != nil {
				v = defaultWaitSeconds
			}
			interval := time.Duration(v) * time.Second

			// Indicate the wait interval
			m.logger.Info(
				logger.EventMonitorStatus,
				fmt.Sprintf(
					"waiting for %s",
					interval,
				),
			)

			// Start the timer
			timer = time.NewTimer(interval)

		case <-m.stopTimerChan:
			if timer != nil {
				timer.Stop()
				timer = nil
			}
		case <-m.closeChan:
			return
		}
	}
}

// New creates a new Monitor instance.
func New(l *logger.Logger, c *conf.Conf, addr, name string) *Monitor {
	var (
		m = &Monitor{
			addr:           addr,
			name:           name,
			logger:         l,
			conf:           c,
			startTimerChan: make(chan any),
			stopTimerChan:  make(chan any),
			closeChan:      make(chan any),
			closedChan:     make(chan any),
		}
		client = nutclient.New(&nutclient.Config{
			Addr:            addr,
			Name:            name,
			ConnectedFn:     m.connectedFn,
			DisconnectedFn:  m.disconnectedFn,
			PowerLostFn:     m.powerLostFn,
			PowerRestoredFn: m.powerRestoredFn,
		})
	)
	m.client = client
	go m.run()
	return m
}

func (m *Monitor) connectedFn() {
	m.logger.Info(
		logger.EventMonitorStatus,
		fmt.Sprintf(
			"connected to %s@%s",
			m.name,
			m.addr,
		),
	)
}

func (m *Monitor) disconnectedFn() {
	m.logger.Warning(
		logger.EventMonitorStatus,
		"disconnected from NUT server",
	)
}

func (m *Monitor) powerLostFn() {

	// Log the powerloss
	m.logger.Error(
		logger.EventMonitorStatus,
		"line power lost; UPS on battery power",
	)

	// Start the action timer
	m.startTimerChan <- nil
}

func (m *Monitor) powerRestoredFn() {

	// Log the power restore event
	m.logger.Info(
		logger.EventMonitorStatus,
		"line power restored",
	)

	// Stop the action timer
	m.stopTimerChan <- nil
}

// Close shuts down the monitor.
func (m *Monitor) Close() {
	m.client.Close()
	close(m.closeChan)
	<-m.closedChan
}
