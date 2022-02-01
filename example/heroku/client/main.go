package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/B3rs/gork/client"
	"github.com/B3rs/gork/example/heroku/client/api"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
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

func main() {

	// open a db connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	manageErr(err)

	c := client.NewClient(db)
	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		},
	))

	// UI routes
	{
		uiHandler, err := getUIHandler(embeddedFiles)
		if err != nil {
			log.Fatal(err)
		}
		e.GET("/*", echo.WrapHandler(uiHandler))
	}

	// API routes
	{
		jobs := api.NewJobsAPI(c)

		v1 := e.Group("/api/v1")
		v1.POST("/jobs/increase", jobs.CreateIncrease)
		v1.GET("/jobs/:id", jobs.Get)
	}

	manageErr(e.Start(":" + os.Getenv("PORT")))
}

func manageErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
