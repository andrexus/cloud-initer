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

	"reflect"
	"strings"

	"github.com/andrexus/cloud-initer/conf"
	"github.com/andrexus/cloud-initer/embedded"
	"github.com/andrexus/cloud-initer/enums"
	"github.com/andrexus/cloud-initer/model"
	"github.com/boltdb/bolt"
	"gopkg.in/go-playground/validator.v9"
)

// API is the data holder for the API
type API struct {
	config *conf.Configuration
	log    *logrus.Entry
	db     *bolt.DB
	echo   *echo.Echo

	// Services used by the API
	instances   model.InstanceService
	environment model.EnvironmentService

	validator CustomValidator
}

type APIListResponse struct {
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"`
	Items    interface{} `json:"items"`
}

type APIResponse struct {
	Status  enums.APIResponseStatus `json:"status"`
	Message string                  `json:"message"`
	Errors  []ErrorResponseItem     `json:"errors,omitempty"`
}

type ErrorResponseItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func NewAPIResponseFromValidationError(errors validator.ValidationErrors) *APIResponse {
	fieldErrors := []ErrorResponseItem{}
	for _, err := range errors {
		fieldName := err.Field()
		message := fmt.Sprintf("Field validation failed on the '%s' validator", fieldName)
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", err.Field())
		case "json":
			message = fmt.Sprintf("%s is not a valid JSON", err.Field())
		case "ip":
			message = fmt.Sprintf("%s is wrong", err.Field())
		case "mac":
			message = fmt.Sprintf("%s is wrong", err.Field())
		}
		if strings.HasPrefix(err.Tag(), "unique") {
			message = fmt.Sprintf("%s '%s' already exists", err.Field(), err.Value())
		}
		fieldErrors = append(fieldErrors, ErrorResponseItem{Field: err.Field(), Message: message})
	}
	response := &APIResponse{Status: enums.Error, Message: "Field validation error", Errors: fieldErrors}
	return response
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
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

	apiValidator := createValidator()
	api.environment = model.NewEnvironmentService(model.NewEnvironmentRepository(db), apiValidator.validator)
	api.instances = model.NewInstanceService(model.NewInstanceRepository(db), api.environment, apiValidator.validator)

	// add the endpoints
	e := echo.New()
	e.HideBanner = true
	e.Validator = apiValidator
	//e.Use(api.logRequest)

	g := e.Group("/api/v1")

	// Instances
	g.GET("/instances", api.InstanceList)
	g.POST("/instances", api.InstanceCreate)
	g.GET("/instances/:id", api.InstanceGet)
	g.PUT("/instances/:id", api.InstanceUpdate)
	g.DELETE("/instances/:id", api.InstanceDelete)
	g.GET("/instances/:id/preview", api.InstancePreview)

	// Environment
	g.GET("/environment", api.EnvironmentGet)
	g.PUT("/environment", api.EnvironmentUpdate)

	// cloud-init
	e.GET("/user-data", api.UserData, api.logRequest, api.injectInstanceByIp)
	e.GET("/meta-data", api.MetaData, api.logRequest, api.injectInstanceByIp)

	e.GET("/*", api.serveVirtualFS, api.frontend404Fallback)

	api.echo = e

	return api
}

func createValidator() *CustomValidator {
	v := CustomValidator{validator: validator.New()}
	v.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &v
}

func (api *API) serveVirtualFS(ctx echo.Context) error {
	w, r := ctx.Response(), ctx.Request()
	fileSystem := embedded.FS(false)
	_, err := fileSystem.Open(r.URL.Path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	fileServer := http.FileServer(fileSystem)
	fileServer.ServeHTTP(w, r)
	return nil
}

func (api *API) frontend404Fallback(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			e, ok := err.(*echo.HTTPError)
			if ok && e.Code == http.StatusNotFound {
				fileSystem := embedded.FS(false)
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
