package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/nathan-osman/nutsvc/conf"
	"github.com/nathan-osman/nutsvc/logger"
	"github.com/nathan-osman/nutsvc/server"
	"github.com/nathan-osman/nutsvc/service"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	serviceName = "nutsvc"
	displayName = "NUT Service"
	description = "Monitor a NUT endpoint for changes"
)

func installService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	p, err := os.Executable()
	if err != nil {
		return err
	}
	s, err := m.CreateService(
		serviceName,
		p,
		mgr.Config{
			StartType:   mgr.StartAutomatic,
			DisplayName: displayName,
			Description: description,
		},
	)
	if err != nil {
		return err
	}
	defer s.Close()
	return eventlog.InstallAsEventCreate(
		serviceName,
		eventlog.Error|eventlog.Warning|eventlog.Info,
	)
}

func serviceCommand(command string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()
	switch command {
	case "remove":
		return s.Delete()
	case "start":
		return s.Start()
	case "stop":
		_, err := s.Control(svc.Stop)
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "nutsvc",
		Usage: "Windows service for monitoring a NUT endpoint",
		Commands: []*cli.Command{
			{
				Name:  "install",
				Usage: "install the application as a service",
				Action: func(c *cli.Context) error {
					return installService()
				},
			},
			{
				Name:  "remove",
				Usage: "remove the service",
				Action: func(c *cli.Context) error {
					return serviceCommand("remove")
				},
			},
			{
				Name:  "start",
				Usage: "start the service",
				Action: func(c *cli.Context) error {
					return serviceCommand("start")
				},
			},
			{
				Name:  "stop",
				Usage: "stop the service",
				Action: func(c *cli.Context) error {
					return serviceCommand("stop")
				},
			},
		},
		Action: func(*cli.Context) error {

			// Application cannot be run interactively
			if i, err := svc.IsWindowsService(); err != nil {
				return err
			} else if !i {
				return errors.New("nutsvc must be run as a Windows service")
			}

			// Create the logger
			l, err := logger.New(serviceName)
			if err != nil {
				return err
			}
			defer l.Close()

			// Create the conf instance
			c, err := conf.New()
			if err != nil {
				return err
			}
			defer c.Close()

			// Create the service
			sInstance := service.New(l)

			// Create the server
			srv := server.New(l)
			defer srv.Close()

			// Run the service
			return svc.Run(serviceName, sInstance)
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
	}
}
