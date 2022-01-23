package web

import (
	"context"
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

type Server struct {
	e *echo.Echo
}

func (s *Server) Start(db *sql.DB, addr string) error {

	s.e = echo.New()

	s.e.HideBanner = true
	s.e.Use(middleware.CORSWithConfig(
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
		s.e.GET("/*", echo.WrapHandler(uiHandler))
		s.e.GET("/jobs/*", echo.WrapHandler(uiHandler))
	}

	// API routes
	{
		jobs := api.NewJobsAPI(db)

		v1 := s.e.Group("/api/v1")
		v1.POST("/jobs/:id/retry", jobs.RetryHandler)
		v1.GET("/jobs/:id", jobs.GetHandler)
		v1.DELETE("/jobs/:id", jobs.CancelHandler)
		v1.GET("/jobs", jobs.ListHandler)
	}

	return s.e.Start(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
