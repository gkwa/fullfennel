package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type App struct {
	statusCheckers []EC2StatusChecker
}

func NewApp(checkers []EC2StatusChecker) *App {
	return &App{statusCheckers: checkers}
}

func (app *App) handleWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		var wg sync.WaitGroup
		for i, checker := range app.statusCheckers {
			wg.Add(1)
			go func(id int, c EC2StatusChecker) {
				defer wg.Done()

				logDir := "log"
				if err := os.MkdirAll(logDir, 0o755); err != nil {
					fmt.Printf("Error creating log directory: %v\n", err)
					return
				}

				logFile, err := os.OpenFile(filepath.Join(logDir, fmt.Sprintf("instance-%02d.log", id)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
				if err != nil {
					fmt.Printf("Error opening log file for instance %d: %v\n", id, err)
					return
				}
				defer logFile.Close()

				logger := log.New(logFile, "", log.LstdFlags)

				var lastStatus InstanceStatus
				for {
					status, err := c.GetEC2Status()
					if err != nil {
						fmt.Printf("Error getting EC2 status for instance %d: %v\n", id, err)
						time.Sleep(1 * time.Second)
						continue
					}

					if status.State != lastStatus.State {
						logger.Printf("State changed from %s to %s", lastStatus.State, status.State)
						lastStatus = status
					}

					err = websocket.JSON.Send(ws, status)
					if err != nil {
						fmt.Printf("Error sending status for instance %d: %v\n", id, err)
						return
					}
					fmt.Printf("Sent status for instance %d: %s\n", id, status.State)

					time.Sleep(1 * time.Second)
				}
			}(i, checker)
		}
		wg.Wait()
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
