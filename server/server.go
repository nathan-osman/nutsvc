package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nathan-osman/nutsvc/logger"
)

// Port is used for the HTTP server.
const Port = 9615

// Server provides the API for interacting with the service.
type Server struct {
	logger *logger.Logger
	server http.Server
}

// New creates a new server instance.
func New(l *logger.Logger) (*Server, error) {
	var (
		m = http.NewServeMux()
		s = &Server{
			logger: l,
			server: http.Server{
				Addr:    fmt.Sprintf(":%d", Port),
				Handler: m,
			},
		}
	)
	go func() {
		defer s.logger.Info(logger.EventServerStatus, "server stopped")
		s.logger.Info(logger.EventServerStatus, "server started")
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(logger.EventServerStatus, err.Error())
		}
	}()
	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Shutdown(context.Background())
}
