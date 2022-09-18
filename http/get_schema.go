package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/d0kur0/webm-api/util"
)

type simplifySchemaItem struct {
	Vendor string   `json:"vendor"`
	Boards []string `json:"board"`
}

func getSchema(w http.ResponseWriter, _ *http.Request) {
	schema := util.ParseSchema()

	var simplifySchema []simplifySchemaItem
	for _, schemaEl := range schema {
		var boards []string
		for _, board := range schemaEl.Boards {
			boards = append(boards, board.Name)
		}
		simplifySchema = append(simplifySchema, simplifySchemaItem{Vendor: schemaEl.Vendor.VendorName(), Boards: boards})
	}

	schemaAsBytes, err := json.Marshal(simplifySchema)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(w, string(schemaAsBytes))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
