//go:build windows

package main

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"time"
)

var elog debug.Log

type exampleService struct {
	httpServer *HttpServer
}

func (m *exampleService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	elog.Info(1, "Starting HTTP server")
	go func() {
		if err := m.httpServer.Start(); err != nil {
			elog.Error(903, fmt.Sprintf("Failed to start server: %s", err))
		}
	}()

	tick := time.Tick(1 * time.Hour)
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case <-tick:
			elog.Info(1, "periodic logic")

		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus

			case svc.Stop, svc.Shutdown:
				elog.Info(2, "Stopping HTTP server")
				if err := m.httpServer.Stop(); err != nil {
					elog.Error(902, fmt.Sprintf("Failed to stop server: %s", err))
				}
				break loop

			case svc.Pause:
				if err := m.httpServer.Stop(); err != nil {
					elog.Error(902, fmt.Sprintf("Failed to stop server: %s", err))
				}
				elog.Info(2, "Stopping HTTP server")
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}

			case svc.Continue:
				elog.Info(1, "Starting HTTP server")
				go func() {
					if err := m.httpServer.Start(); err != nil {
						elog.Error(903, fmt.Sprintf("Failed to start server: %s", err))
					}
				}()
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

			default:
				elog.Error(901, fmt.Sprintf("Unexpected control request #%d", c))
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}

	err = run(name, &exampleService{
		httpServer: NewHttpServer(2025),
	})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}

func main() {
	runService("poc-go-windows-service", false)
}
