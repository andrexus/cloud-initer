package api

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"context"

	"net/http"

	"bytes"
	"io"

	"github.com/andrexus/cloud-initer/conf"
	"github.com/andrexus/cloud-initer/generated"
	"github.com/andrexus/cloud-initer/model"
	"github.com/boltdb/bolt"
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
	//e.Use(api.logRequest)

	g := e.Group("/api/v1")

	// Instances
	g.GET("/instances", api.InstanceList)
	g.POST("/instances", api.InstanceCreate)
	g.GET("/instances/:id", api.InstanceGet)
	g.PUT("/instances/:id", api.InstanceUpdate)
	g.DELETE("/instances/:id", api.InstanceDelete)

	// cloud-init
	e.GET("/user-data", api.UserData, api.logRequest, api.injectInstanceByIp)
	e.GET("/meta-data", api.MetaData, api.logRequest, api.injectInstanceByIp)

	e.GET("/*", api.serveVirtualFS, api.angularRouterFallback)

	api.echo = e

	return api
}

func (api *API) serveVirtualFS(ctx echo.Context) error {
	w, r := ctx.Response(), ctx.Request()
	fileSystem := generated.FS(false)
	_, err := fileSystem.Open(r.URL.Path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	fileServer := http.FileServer(fileSystem)
	fileServer.ServeHTTP(w, r)
	return nil
}

func (api *API) angularRouterFallback(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			e, ok := err.(*echo.HTTPError)
			if ok && e.Code == http.StatusNotFound {
				fileSystem := generated.FS(false)
				f, _ := fileSystem.Open("/index.html")
				buf := bytes.NewBuffer(nil)
				io.Copy(buf, f)
				f.Close()
				c.HTML(http.StatusOK, string(buf.Bytes()))
			}
		}
		return err
	}
}

func (api *API) logRequest(f echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		logger := api.log.WithFields(logrus.Fields{
			"method": req.Method,
			"path":   req.URL.Path,
		})
		ctx.Set(loggerKey, logger)

		logger.WithFields(logrus.Fields{
			"user_agent": req.UserAgent(),
			"ip_address": ctx.RealIP(),
		}).Info("Request")

		err := f(ctx)
		if err != nil {
			ctx.Error(err)
		}
		return err
	}
}
