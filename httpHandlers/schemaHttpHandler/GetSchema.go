package schemaHttpHandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/d0kur0/webm-api/tasks/grabberTask"
)

func GetSchema(w http.ResponseWriter, r *http.Request) {
	var grabberSchema = grabberTask.GetGrabberSchema()
	var responseData = make(responseSchema, len(grabberSchema))

	for _, schema := range grabberSchema {
		var boards []responseBoard
		for _, board := range schema.Boards {
			boards = append(boards, responseBoard{board.Name, board.Description})
		}

		responseData[schema.Vendor.VendorName()] = boards
	}

	jsonBytes, err := json.Marshal(responseData)
	if err != nil {
		log.Println("Marshaling output failed:", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(jsonBytes); err != nil {
		log.Println("Writing response failed:", err)
	}
}
