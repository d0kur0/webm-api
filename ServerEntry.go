package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d0kur0/webm-api/httpHandlers/filesHttpHandler"
	"github.com/d0kur0/webm-api/httpHandlers/schemaHttpHandler"
	"github.com/d0kur0/webm-api/tasks/grabberTask"

	"github.com/gorilla/mux"
)

const port = "3500"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	go grabberTask.Start()

	router := mux.NewRouter()
	router.Use(accessControlMiddleware)
	router.HandleFunc("/schema/get", schemaHttpHandler.GetSchema).Methods("GET")
	router.HandleFunc("/files/getByStruct", filesHttpHandler.GetFilesByStruct).Methods("POST", "OPTIONS")
	router.HandleFunc("/files/getAll", filesHttpHandler.GetAll).Methods("GET")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Starting server failed: ", err)
	}
	log.Println("Server started at port:", port)
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
