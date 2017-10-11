package api

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"context"

	"github.com/andrexus/cloud-initer/conf"
	"github.com/andrexus/cloud-initer/model"
	"github.com/boltdb/bolt"
	"github.com/pborman/uuid"
)

// API is the data holder for the API
type API struct {
	config *conf.Configuration
	log    *logrus.Entry
	db     *bolt.DB
	echo   *echo.Echo

	// Services used by the API
	instances model.InstanceService
}

type APIListResponse struct {
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"`
	Items    interface{} `json:"items"`
}

type APIResponse struct {
	Message string `json:"message"`
}

// Start will start the API on the specified port
func (api *API) Start() error {
	return api.echo.Start(fmt.Sprintf(":%d", api.config.API.Port))
}

// Stop will shutdown the engine internally
func (api *API) Stop() error {
	logrus.Info("Stopping API server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return api.echo.Shutdown(ctx)
}

// NewAPI will create an api instance that is ready to start
func NewAPI(config *conf.Configuration, db *bolt.DB) *API {
	api := &API{
		config: config,
		log:    logrus.WithField("component", "api"),
		db:     db,
	}

	api.instances = model.NewInstanceService(model.NewInstanceRepository(db))

	// add the endpoints
	e := echo.New()
	e.HideBanner = true
	//e.Use(api.setupRequest)

	g := e.Group("/api/v1")

	// Instances
	g.GET("/instances", api.InstanceList)
	g.POST("/instances", api.InstanceCreate)
	g.GET("/instances/:id", api.InstanceGet)
	g.PUT("/instances/:id", api.InstanceUpdate)
	g.DELETE("/instances/:id", api.InstanceDelete)

	api.echo = e

	return api
}

func (api *API) setupRequest(f echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		logger := api.log.WithFields(logrus.Fields{
			"method":     req.Method,
			"path":       req.URL.Path,
			"request_id": uuid.NewRandom().String(),
		})
		ctx.Set(loggerKey, logger)

		startTime := time.Now()
		defer func() {
			rsp := ctx.Response()
			logger.WithFields(logrus.Fields{
				"status_code":  rsp.Status,
				"runtime_nano": time.Since(startTime).Nanoseconds(),
			}).Info("Finished request")
		}()

		logger.WithFields(logrus.Fields{
			"user_agent":     req.UserAgent(),
			"content_length": req.ContentLength,
		}).Info("Starting request")

		// we have to do this b/c if not the final error handler will not
		// in the chain of middleware. It will be called after meaning that the
		// response won't be set properly.
		err := f(ctx)
		if err != nil {
			ctx.Error(err)
		}
		return err
	}
}
