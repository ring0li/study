package common

import (
	//"gitflow/config"
	"github.com/labstack/echo/v4"
	//"go.uber.org/zap"
	"net/http"
)

type RetStruct struct {
	ErrorCode int         `json:"errorCode"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

func Succ(data interface{}) *RetStruct {
	//
	//config.Logger.Info("Success", zap.Any("data", data))
	//
	return &RetStruct{0, "", data}
}

func Fail(code int, message string) *RetStruct {
	//
	//config.Logger.Info("Failure", zap.Int("code", code), zap.String("message", message))
	//
	return &RetStruct{code, message, []string{}}
}

func Output(c echo.Context, ret *RetStruct) error {

	return c.JSON(http.StatusOK, ret)
}
