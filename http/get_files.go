package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/d0kur0/webm-api/worker"
)

func getFiles(c echo.Context) error {
	return c.JSON(http.StatusOK, worker.GrabbingOutPut)
}
