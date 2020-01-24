package filesHttpHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/d0kur0/webm-api/tasks/grabberTask"

	"github.com/ztrue/tracerr"
)

func GetFilesByStruct(w http.ResponseWriter, r *http.Request) {
	var requestSchema requestSchema
	var requestBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(requestBody, &requestSchema); err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var output = grabberTask.GetOutput()
	for vendor, boards := range output.Vendors {
		if _, exists := requestSchema.Vendors[vendor]; !exists {
			delete(output.Vendors, vendor)
			continue
		}

		for board := range boards {
			exists := false
			for _, requestBoard := range requestSchema.Vendors[vendor] {
				if requestBoard == board {
					exists = true
					break
				}
			}

			if !exists {
				delete(output.Vendors[vendor], board)
			}
		}
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		tracerr.PrintSourceColor(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(jsonBytes); err != nil {
		tracerr.PrintSourceColor(tracerr.New("Write response error"))
	}
}
