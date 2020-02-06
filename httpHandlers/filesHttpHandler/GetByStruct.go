package filesHttpHandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/d0kur0/webm-grabber/sources/types"

	"github.com/d0kur0/webm-api/tasks/grabberTask"
)

func GetFilesByStruct(w http.ResponseWriter, r *http.Request) {
	var requestSchema requestSchema
	var requestBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Reading response body failed:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(requestBody, &requestSchema); err != nil {
		log.Println("Unmarshal request body failed:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var responseOutput = types.Output{}
	responseOutput.Vendors = make(types.OutputVendors, len(requestSchema.Vendors))

	for vendor, boards := range grabberTask.GetOutput().Vendors {
		if _, exists := requestSchema.Vendors[vendor]; exists {
			var outputBoards []types.OutputBoard

			for _, board := range boards {
				exists := false
				for _, requestBoard := range requestSchema.Vendors[vendor] {
					if requestBoard == board.Name {
						exists = true
						break
					}
				}

				if exists {
					outputBoards = append(outputBoards, types.OutputBoard{
						Name:        board.Name,
						Description: board.Description,
						Threads:     board.Threads,
					})
				}
			}

			responseOutput.Vendors[vendor] = outputBoards
		}
	}

	jsonBytes, err := json.Marshal(responseOutput)
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
