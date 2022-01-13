package web

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/B3rs/gork/client"
	"github.com/B3rs/gork/jobs"
	echo "github.com/labstack/echo/v4"
)

func Start(db *sql.DB, addr string) error {
	e := echo.New()

	jobsRoutes := NewJobsRoute(db)

	e.POST("/api/v1/jobs/:id/retry", jobsRoutes.retry)
	e.GET("/api/v1/jobs/:id", jobsRoutes.get)
	e.DELETE("/api/v1/jobs/:id", jobsRoutes.cancel)
	e.GET("/api/v1/jobs", jobsRoutes.list)

	return e.Start(addr)
}

func NewJobsRoute(db *sql.DB) JobsRoute {
	return JobsRoute{client: client.NewDBClient(db)}
}

type JobsRoute struct {
	client client.Client
}

type jobsList struct {
	Jobs []*jobs.Job `json:"jobs"`
}

func (jr JobsRoute) list(c echo.Context) error {

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	page, _ := strconv.Atoi(c.QueryParam("page"))

	jobs, err := jr.client.GetAll(c.Request().Context(), page, limit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jobsList{jobs})
}

func (jr JobsRoute) get(c echo.Context) error {

	job, err := jr.client.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, job)
}

func (jr JobsRoute) retry(c echo.Context) error {

	id := c.Param("id")

	err := jr.client.ForceRetry(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (jr JobsRoute) cancel(c echo.Context) error {

	id := c.Param("id")

	err := jr.client.Cancel(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
