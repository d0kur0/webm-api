package httpActions

import (
	"encoding/json"
	"net/http"
	"webmApi/filesDaemon"

	"github.com/ztrue/tracerr"
)

func GetSchema(response http.ResponseWriter, request *http.Request) {
	var grabberSchema = filesDaemon.GetGrabberSchema()
	var responseData = make(map[string][]string, len(grabberSchema))

	for _, schema := range grabberSchema {
		var boards []string
		for _, board := range schema.Boards {
			boards = append(boards, board.String())
		}

		responseData[schema.Vendor.VendorName()] = boards
	}

	jsonBytes, err := json.Marshal(responseData)
	if err != nil {
		tracerr.PrintSourceColor(tracerr.New("JSON marshal error"))
		http.Error(response, "Server Error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")

	if _, err = response.Write(jsonBytes); err != nil {
		tracerr.PrintSourceColor(tracerr.New("Write response error"))
	}
}
