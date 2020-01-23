package main

import (
	"fmt"
	"net/http"
	"webmApi/filesDaemon"
	"webmApi/httpActions"

	"github.com/gorilla/mux"
	"github.com/ztrue/tracerr"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/getSchema", httpActions.GetSchema)
	r.HandleFunc("/getFiles", httpActions.GetFiles)

	filesDaemon.Start()

	if err := http.ListenAndServe(":3500", r); err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
	}
}
