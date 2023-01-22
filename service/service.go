package service

import (
	"github.com/nathan-osman/nutsvc/logger"
	"golang.org/x/sys/windows/svc"
)

// Service monitors for UPS state changes.
type Service struct {
	logger *logger.Logger
}

// New creates a new Service instance.
func New(l *logger.Logger) *Service {
	return &Service{
		logger: l,
	}
}

func (s *Service) Execute(args []string, chChan <-chan svc.ChangeRequest, stChan chan<- svc.Status) (bool, uint32) {
	defer s.logger.Info(logger.EventServiceStatus, "service stopped")
	s.logger.Info(logger.EventServiceStatus, "service started")

	// Indicate that the service has been started
	stChan <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

	// Respond to service requests
	for c := range chChan {
		switch c.Cmd {
		case svc.Interrogate:
			stChan <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			stChan <- svc.Status{
				State: svc.StopPending,
			}
			return false, 0
		}
	}

	// This line should never be reached
	return false, 0
}
