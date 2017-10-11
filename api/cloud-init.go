package api

import (
	"net/http"

	"github.com/andrexus/cloud-initer/model"
	"github.com/labstack/echo"
)

const instanceKey = "request.instance"

func (api *API) UserData(ctx echo.Context) error {
	item := ctx.Get(instanceKey)
	return ctx.String(http.StatusOK, item.(*model.Instance).UserData)
}

func (api *API) MetaData(ctx echo.Context) error {
	item := ctx.Get(instanceKey)
	return ctx.String(http.StatusOK, item.(*model.Instance).MetaData)
}

func (api *API) injectInstanceByIp(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ip := ctx.RealIP()
		item, _ := api.instances.FindByIP(ip)
		if item == nil {
			response := &APIResponse{Message: "instance not found"}
			return ctx.JSON(http.StatusNotFound, response)
		}
		ctx.Set(instanceKey, item)
		err := next(ctx)
		if err != nil {
			ctx.Error(err)
		}
		return err
	}
}
