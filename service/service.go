package service

import (
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	// Name identifies the service within the SCM.
	Name = "nutsvc"

	EventServiceStatus = 1
)

// Service monitors for UPS state changes.
type Service struct {
	log      *eventlog.Log
	stopChan chan any
}

// New creates a new Service instance.
func New() (*Service, error) {
	e, err := eventlog.Open(Name)
	if err != nil {
		return nil, err
	}
	s := &Service{
		log:      e,
		stopChan: make(chan any),
	}
	return s, nil
}

func (s *Service) Execute(args []string, chChan <-chan svc.ChangeRequest, stChan chan<- svc.Status) (bool, uint32) {
	s.log.Info(EventServiceStatus, "event loop started")
	defer s.log.Info(EventServiceStatus, "event loop stopped")

	// Indicate that the service has been started
	stChan <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

	// Indicate that the service is stopping when exiting
	defer func() {
		stChan <- svc.Status{
			State: svc.StopPending,
		}
	}()

	for {
		select {
		case <-s.stopChan:
			return false, 0
		case c := <-chChan:
			switch c.Cmd {
			case svc.Interrogate:
				stChan <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				return false, 0
			}
		}
	}
}

// Close shuts down the service.
func (s *Service) Close() {
	close(s.stopChan)
}
