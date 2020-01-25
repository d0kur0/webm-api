package filesHttpHandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/d0kur0/webm-api/tasks/grabberTask"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	var output = grabberTask.GetOutput()

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
