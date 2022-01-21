package web

import (
	"database/sql"
	"embed"
	"io/fs"
	"net/http"

	"github.com/B3rs/gork/web/api"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// content holds our static web server content.
//go:embed ui/build/*
var embeddedFiles embed.FS

func getUIHandler(fsys fs.FS) (http.Handler, error) {
	buildDir, err := fs.Sub(fsys, "ui/build")
	if err != nil {
		return nil, err
	}

	return http.FileServer(http.FS(buildDir)), nil
}

func Start(db *sql.DB, addr string) error {

	e := echo.New()
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		},
	))

	// UI routes
	{
		uiHandler, err := getUIHandler(embeddedFiles)
		if err != nil {
			return err
		}
		e.GET("/*", echo.WrapHandler(uiHandler))
	}

	// API routes
	{
		jobs := api.NewJobsAPI(db)

		v1 := e.Group("/api/v1")
		v1.POST("/jobs/:id/retry", jobs.RetryHandler)
		v1.GET("/jobs/:id", jobs.GetHandler)
		v1.DELETE("/jobs/:id", jobs.CancelHandler)
		v1.GET("/jobs", jobs.ListHandler)
	}

	return e.Start(addr)
}
