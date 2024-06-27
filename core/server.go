package core

import (
	"embed"
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed index.html
var content embed.FS

const mockServerCount = 100

func Run() {
	os.RemoveAll("log")

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var statusCheckers []EC2StatusChecker
	for i := 0; i < mockServerCount; i++ {
		statusChecker := NewMockEC2StatusChecker(
			fmt.Sprintf("mock-instance-%02d", i),
			rng,
		)
		statusCheckers = append(statusCheckers, statusChecker)
	}

	app := NewApp(statusCheckers)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fsys, err := fs.Sub(content, ".")
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", echo.WrapHandler(http.FileServer(http.FS(fsys))))
	e.GET("/ws", app.handleWebSocket)

	e.Logger.Fatal(e.Start(":8080"))
}
