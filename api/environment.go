package api

import (
	"net/http"

	"github.com/andrexus/cloud-initer/model"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func (api *API) EnvironmentGet(ctx echo.Context) error {
	item, err := api.environment.GetEnvironment()

	if err != nil {
		response := &APIResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	return ctx.JSON(http.StatusOK, item)

}

func (api *API) EnvironmentUpdate(ctx echo.Context) error {
	item := new(model.Environment)
	if err := ctx.Bind(item); err != nil {
		response := &APIResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if err := ctx.Validate(item); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIResponseFromValidationError(err.(validator.ValidationErrors)))
	}
	item, err := api.environment.Update(item)
	if err != nil {
		response := &APIResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	return ctx.JSON(http.StatusOK, item)

}
