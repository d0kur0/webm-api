package webm_api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d0kur0/webm-api/httpHandlers/filesHttpHandler"
	"github.com/d0kur0/webm-api/httpHandlers/schemaHttpHandler"
	"github.com/d0kur0/webm-api/tasks/grabberTask"

	"github.com/gorilla/mux"
	"github.com/ztrue/tracerr"
)

const port = "3500"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	go grabberTask.Start()

	r := mux.NewRouter()
	r.HandleFunc("/schema/get", schemaHttpHandler.GetSchema).Methods("GET")
	r.HandleFunc("/files/getByStruct", filesHttpHandler.GetFilesByStruct).Methods("POST")
	r.HandleFunc("/files/getAll", filesHttpHandler.GetAll).Methods("GET")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
	}

	log.Println("Server started at " + port + " port")
}
