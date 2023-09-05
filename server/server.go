package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nathan-osman/nutsvc/conf"
	"github.com/nathan-osman/nutsvc/logger"
	"github.com/nathan-osman/nutsvc/monitor"
)

// Port is used for the HTTP server.
const Port = 9615

// Switch Gin to release mode
func init() {
	gin.SetMode(gin.ReleaseMode)
}

// Server provides the API for interacting with the service.
type Server struct {
	logger  *logger.Logger
	conf    *conf.Conf
	monitor *monitor.Monitor
	server  http.Server
}

func (s *Server) tryReloadConf() error {
	err := func() error {
		if s.monitor != nil {
			s.monitor.Close()
		}
		a, err := s.conf.Get("addr")
		if err != nil {
			return err
		}
		n, err := s.conf.Get("name")
		if err != nil {
			return err
		}
		s.monitor = monitor.New(s.logger, s.conf, a, n)
		return nil
	}()
	if err != nil {
		s.logger.Error(
			logger.EventServerStatus,
			fmt.Sprintf(
				"unable to load config: %s",
				err.Error(),
			),
		)
	}
	return err
}

// New creates a new server instance.
func New(l *logger.Logger, c *conf.Conf) *Server {
	var (
		g = gin.New()
		s = &Server{
			logger: l,
			conf:   c,
			server: http.Server{
				Addr:    fmt.Sprintf(":%d", Port),
				Handler: g,
			},
		}
	)

	// Attempt to reload the config
	s.tryReloadConf()

	// API methods
	api := g.Group("/api")
	api.Use(gin.CustomRecoveryWithWriter(nil, panicToJSONError))
	{
		api.GET("/config", s.apiConfig_GET)
		api.POST("/config", s.apiConfig_POST)
		api.POST("/monitor/reload", s.apiMonitorReload_POST)
	}

	go func() {
		defer s.logger.Info(logger.EventServerStatus, "server stopped")
		s.logger.Info(
			logger.EventServerStatus,
			fmt.Sprintf("server started on port %d", Port),
		)
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(logger.EventServerStatus, err.Error())
		}
	}()
	return s
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Shutdown(context.Background())
}
