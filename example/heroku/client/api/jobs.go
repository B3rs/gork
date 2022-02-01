package api

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/B3rs/gork/client"
	herokujobs "github.com/B3rs/gork/example/heroku/jobs"
	"github.com/B3rs/gork/jobs"
	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type scheduler interface {
	Schedule(ctx context.Context, id string, queueName string, arguments interface{}, options ...client.OptionFunc) error
	Get(ctx context.Context, id string) (jobs.Job, error)
}

func NewJobsAPI(scheduler scheduler) JobsAPI {
	return JobsAPI{scheduler: scheduler}
}

type JobsAPI struct {
	scheduler scheduler
}

type createParams struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
}

func (j JobsAPI) CreateIncrease(c echo.Context) error {

	params := &createParams{}
	err := c.Bind(params)
	if err != nil {
		return err
	}

	if params.ID == "" {
		params.ID = uuid.New().String()
	}
	if params.Number == 0 {
		params.Number = rand.Intn(1000)
	}

	err = j.scheduler.Schedule(c.Request().Context(), params.ID, "increase", herokujobs.IncreaseArgs{IncreaseThis: params.Number})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, params)
}

func (j JobsAPI) Get(c echo.Context) error {

	job, err := j.scheduler.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, job)
}