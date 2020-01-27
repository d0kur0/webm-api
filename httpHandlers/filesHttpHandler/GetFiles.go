package filesHttpHandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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

	var output = grabberTask.GetOutput()
	for vendor, boards := range output.Vendors {
		if _, exists := requestSchema.Vendors[vendor]; !exists {
			delete(output.Vendors, vendor)
			continue
		}

		for boardIndex, board := range boards {
			exists := false
			for _, requestBoard := range requestSchema.Vendors[vendor] {
				if requestBoard == board.Name {
					exists = true
					break
				}
			}

			if !exists {
				output.Vendors[vendor][len(output.Vendors[vendor])-1], output.Vendors[vendor][boardIndex] = output.Vendors[vendor][boardIndex], output.Vendors[vendor][len(output.Vendors[vendor])-1]
				output.Vendors[vendor] = output.Vendors[vendor][:len(output.Vendors[vendor])-1]
			}
		}
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		log.Println("Marshaling output failed:", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if _, err = w.Write(jsonBytes); err != nil {
		log.Println("Writing response failed:", err)
	}
}
