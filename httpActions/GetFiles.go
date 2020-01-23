package httpActions

import (
	"encoding/json"
	"net/http"
	"webmApi/filesDaemon"

	"github.com/ztrue/tracerr"
)

func GetFiles(response http.ResponseWriter, request *http.Request) {
	var files = filesDaemon.GetOutput()

	jsonBytes, err := json.Marshal(files)
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
