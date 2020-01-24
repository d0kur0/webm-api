package filesHttpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/d0kur0/webm-api/tasks/grabberTask"

	"github.com/ztrue/tracerr"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	var output = grabberTask.GetOutput()

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
