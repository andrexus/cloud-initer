package api

import (
	"net/http"

	"github.com/andrexus/cloud-initer/enums"
	"github.com/andrexus/cloud-initer/model"
	"github.com/labstack/echo"
)

const instanceKey = "request.instance"

func (api *API) Preview(ctx echo.Context) error {
	data := new(model.CloudInitData)
	if err := ctx.Bind(data); err != nil {
		response := &APIResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	result, err := api.cloudInit.PreviewCloudInitData(data.UserData, data.MetaData)
	if err != nil {
		response := &APIResponse{Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	return ctx.JSON(http.StatusOK, result)

}

func (api *API) UserData(ctx echo.Context) error {
	item := ctx.Get(instanceKey).(*model.CloudInitData)
	return ctx.String(http.StatusOK, item.UserData)
}

func (api *API) MetaData(ctx echo.Context) error {
	item := ctx.Get(instanceKey).(*model.CloudInitData)
	return ctx.String(http.StatusOK, item.MetaData)
}

func (api *API) injectInstanceByIp(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ip := ctx.RealIP()
		cloudInitData, e := api.cloudInit.GetCloudInitDataForClient(ip, ctx.Request().UserAgent())
		if e != nil {
			response := &APIResponse{Status: enums.Error, Message: e.Error()}
			return ctx.JSON(http.StatusInternalServerError, response)
		}

		ctx.Set(instanceKey, cloudInitData)
		err := next(ctx)
		if err != nil {
			ctx.Error(err)
		}
		return err
	}
}
