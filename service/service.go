package service

import (
	"github.com/nathan-osman/nutsvc/logger"
	"golang.org/x/sys/windows/svc"
)

// Service monitors for UPS state changes.
type Service struct {
	logger   *logger.Logger
	stopChan chan any
}

// New creates a new Service instance.
func New(l *logger.Logger) (*Service, error) {
	s := &Service{
		logger:   l,
		stopChan: make(chan any),
	}
	return s, nil
}

func (s *Service) Execute(args []string, chChan <-chan svc.ChangeRequest, stChan chan<- svc.Status) (bool, uint32) {
	defer s.logger.Info(logger.EventServiceStatus, "service stopped")
	s.logger.Info(logger.EventServiceStatus, "service started")

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
