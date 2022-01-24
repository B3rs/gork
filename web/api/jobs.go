package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/B3rs/gork/db"
	"github.com/B3rs/gork/jobs"
	echo "github.com/labstack/echo/v4"
)

type dbStore interface {
	Search(ctx context.Context, limit int, offset int, search string) ([]jobs.Job, error)
	Get(ctx context.Context, id string) (jobs.Job, error)
	Update(ctx context.Context, job jobs.Job) error
	Create(ctx context.Context, job jobs.Job) error
	Deschedule(ctx context.Context, id string) error
	ScheduleNow(ctx context.Context, id string) error
}

func NewJobsAPI(database *sql.DB) JobsAPI {
	return JobsAPI{db: db.NewStore(database)}
}

type JobsAPI struct {
	db dbStore
}

type jobsList struct {
	Jobs []jobs.Job `json:"jobs"`
}

func (j JobsAPI) ListHandler(c echo.Context) error {

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 50
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	search := c.QueryParam("q")

	jobs, err := j.db.Search(c.Request().Context(), page-1, limit, search)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jobsList{jobs})
}

func (j JobsAPI) GetHandler(c echo.Context) error {

	job, err := j.db.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, job)
}

func (j JobsAPI) RetryHandler(c echo.Context) error {

	id := c.Param("id")

	err := j.db.ScheduleNow(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (j JobsAPI) CancelHandler(c echo.Context) error {

	id := c.Param("id")

	err := j.db.Deschedule(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
