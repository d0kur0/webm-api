package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/d0kur0/webm-api/worker"
)

func getFiles(w http.ResponseWriter, _ *http.Request) {
	outputAsBytes, err := json.Marshal(worker.GrabbingOutPut)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(w, string(outputAsBytes))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
