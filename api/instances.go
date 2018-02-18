package api

import (
	"net/http"

	"github.com/andrexus/cloud-initer/model"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func (api *API) InstanceList(ctx echo.Context) error {
	var err error

	items, err := api.instances.FindAll()
	if err != nil {
		response := &MessageResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	response := &ListResponse{Page: 1, PageSize: len(items), Total: len(items), Items: items}
	return ctx.JSON(http.StatusOK, response)
}

func (api *API) InstanceCreate(ctx echo.Context) error {
	item := new(model.Instance)
	if err := ctx.Bind(item); err != nil {
		response := &MessageResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if err := ctx.Validate(item); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIResponseFromValidationError(err.(validator.ValidationErrors)))
	}
	item, err := api.instances.Create(item)
	if err != nil {
		response := &MessageResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	return ctx.JSON(http.StatusCreated, item)

}

func (api *API) InstanceGet(ctx echo.Context) error {
	id := ctx.Param("id")
	item, err := api.instances.FindOne(id)

	if err != nil {
		response := &MessageResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if item == nil {
		response := &MessageResponse{Message: "instance not found"}
		return ctx.JSON(http.StatusNotFound, response)
	}
	return ctx.JSON(http.StatusOK, item)

}

func (api *API) InstanceUpdate(ctx echo.Context) error {
	id := ctx.Param("id")
	newItem := new(model.Instance)
	if err := ctx.Bind(newItem); err != nil {
		response := &MessageResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if err := ctx.Validate(newItem); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIResponseFromValidationError(err.(validator.ValidationErrors)))
	}
	item, err := api.instances.Update(id, newItem)
	if err != nil {
		response := &MessageResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	return ctx.JSON(http.StatusOK, item)

}

func (api *API) InstanceDelete(ctx echo.Context) error {
	id := ctx.Param("id")
	err := api.instances.Delete(id)
	if err != nil {
		response := &MessageResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	response := &MessageResponse{Message: "instance deleted"}
	return ctx.JSON(http.StatusOK, response)
}
