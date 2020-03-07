package main

import (
	"fmt"
	"github.com/d0kur0/webm-api/helpers/settingsManager"
	"github.com/d0kur0/webm-api/httpHandlers/filesHttpHandler"
	"github.com/d0kur0/webm-api/httpHandlers/schemaHttpHandler"
	"github.com/d0kur0/webm-api/tasks/grabberTask"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	config, err := settingsManager.Parse()
	if err != nil {
		log.Println(err)
		return
	}

	go grabberTask.Start(config.UpdateInterval)

	router := mux.NewRouter()
	router.Use(accessControlMiddleware)
	router.HandleFunc("/schema/get", schemaHttpHandler.GetSchema).Methods("GET")
	router.HandleFunc("/files/getByStruct", filesHttpHandler.GetFilesByStruct).Methods("POST", "OPTIONS")
	router.HandleFunc("/files/getAll", filesHttpHandler.GetAll).Methods("GET")

	log.Println("Server started at port:", config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handlers.CompressHandler(router)); err != nil {
		log.Println("Starting server failed: ", err)
	}
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
