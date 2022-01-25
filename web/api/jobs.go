package api

import (
	"net/http"
	"strconv"

	"github.com/B3rs/gork/db"
	"github.com/B3rs/gork/jobs"
	echo "github.com/labstack/echo/v4"
)

func NewJobsAPI(store db.JobsStore) JobsAPI {
	return JobsAPI{store: store}
}

type JobsAPI struct {
	store db.JobsStore
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

	jobs, err := j.store.Search(c.Request().Context(), page-1, limit, search)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jobsList{jobs})
}

func (j JobsAPI) GetHandler(c echo.Context) error {

	job, err := j.store.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, job)
}

func (j JobsAPI) RetryHandler(c echo.Context) error {

	id := c.Param("id")

	err := j.store.ScheduleNow(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (j JobsAPI) CancelHandler(c echo.Context) error {

	id := c.Param("id")

	err := j.store.Deschedule(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
