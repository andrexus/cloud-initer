package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

const (
	loggerKey = "app.logger"
)

func getLogger(ctx echo.Context) *logrus.Entry {
	obj := ctx.Get(loggerKey)
	if obj == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return obj.(*logrus.Entry)
}
