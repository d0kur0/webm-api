package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/d0kur0/webm-api/util"
)

type simplifySchemaItem struct {
	Vendor string   `json:"vendor"`
	Boards []string `json:"boards"`
}

func getSchema(c echo.Context) error {
	schema := util.ParseSchema()

	var simplifySchema []simplifySchemaItem
	for _, schemaEl := range schema {
		var boards []string
		for _, board := range schemaEl.Boards {
			boards = append(boards, board.Name)
		}
		simplifySchema = append(simplifySchema, simplifySchemaItem{Vendor: schemaEl.Vendor.VendorName(), Boards: boards})
	}

	return c.JSON(http.StatusOK, simplifySchema)
}
