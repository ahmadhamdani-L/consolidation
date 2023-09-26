package middleware

import (
	"fmt"
	"mcash-finance-console-core/configs"
	"mcash-finance-console-core/pkg/util/validator"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo) {
	NAME := fmt.Sprintf("%s-%s", configs.App().Name(), configs.App().Env())

	e.Use(Context)
	e.Use(
		echoMiddleware.Recover(),
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodDelete},
		}),
		echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format:           fmt.Sprintf("\n%s | ${host} | ${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri}", NAME),
			CustomTimeFormat: "2006/01/02 15:04:05",
			Output:           os.Stdout,
		}),
	)
	e.HTTPErrorHandler = ErrorHandler
	e.Validator = &validator.CustomValidator{Validator: validator.NewValidator()}
}
