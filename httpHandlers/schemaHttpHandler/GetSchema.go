package schemaHttpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/d0kur0/webm-api/tasks/grabberTask"

	"github.com/ztrue/tracerr"
)

func GetSchema(w http.ResponseWriter, r *http.Request) {
	var grabberSchema = grabberTask.GetGrabberSchema()
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
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if _, err = w.Write(jsonBytes); err != nil {
		tracerr.PrintSourceColor(tracerr.New("Write response error"))
	}
}
