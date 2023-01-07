package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/d0kur0/webm-api/worker"

	"github.com/d0kur0/webm-grabber/types"
)

type getFilesConditionRequest map[string][]string

func getFilesByCondition(c echo.Context) error {
	var condition getFilesConditionRequest
	if err := c.Bind(&condition); err != nil {
		return err
	}

	var filteredFiles types.Output
	for _, outItem := range worker.GrabbingOutPut {
		boards, isDesiredVendor := condition[outItem.VendorName]
		if !isDesiredVendor {
			continue
		}

		isDesiredBoard := false
		for _, board := range boards {
			if board == outItem.BoardName {
				isDesiredBoard = true
				break
			}
		}

		if !isDesiredBoard {
			continue
		}

		filteredFiles = append(filteredFiles, outItem)
	}

	return c.JSON(http.StatusOK, filteredFiles)
}
