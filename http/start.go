package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"

	"github.com/spf13/viper"
)

func Start() error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	e.GET("/schema", getSchema)
	e.GET("/files", getFiles)
	e.POST("/filesByCondition", getFilesByCondition)

	return e.Start(fmt.Sprintf(":%d", viper.GetInt("port")))
}
