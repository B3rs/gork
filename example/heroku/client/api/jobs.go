package api

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/B3rs/gork/client"
	herokujobs "github.com/B3rs/gork/example/heroku/jobs"
	"github.com/B3rs/gork/jobs"
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
	ID          string    `json:"id"`
	Queue       string    `json:"queue"`
	Number      int       `json:"number,omitempty"`
	String      string    `json:"string,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

func (j JobsAPI) Create(c echo.Context) error {

	params := new(createParams)
	if err := c.Bind(params); err != nil {
		return err
	}

	options := []client.OptionFunc{}
	if params.ScheduledAt.IsZero() {
		params.ScheduledAt = time.Now().Add(time.Duration(rand.Intn(10)) * time.Second)
	}
	options = append(options, client.WithScheduleTime(params.ScheduledAt))

	var args interface{}
	if params.Queue == "increase" {
		args = herokujobs.IncreaseArgs{IncreaseThis: params.Number}
	}
	if params.Queue == "lowerize" {
		args = herokujobs.LowerizeArgs{LowerizeThis: params.String}
	}

	if err := j.scheduler.Schedule(
		c.Request().Context(),
		params.ID,
		params.Queue,
		args,
		options...,
	); err != nil {
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
