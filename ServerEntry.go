package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/d0kur0/webm-api/httpHandlers/filesHttpHandler"
	"github.com/d0kur0/webm-api/httpHandlers/schemaHttpHandler"
	"github.com/d0kur0/webm-api/tasks/grabberTask"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Port           int
	UpdateInterval uint64
}

func parseServerConfig() (config ServerConfig, err error) {
	const configFilePath = "serverConfig.json"
	const defaultPort = 3500
	const defaultUpdateInterval = 10

	config = ServerConfig{
		Port:           defaultPort,
		UpdateInterval: defaultUpdateInterval,
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		jsonBytes, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return config, errors.Wrap(err, fmt.Sprintf("Marshaling config error (%s)", configFilePath))
		}

		if err := ioutil.WriteFile(configFilePath, jsonBytes, 0644); err != nil {
			return config, errors.Wrap(err, fmt.Sprintf("Write in config data file error (%s)", configFilePath))
		}
	} else {
		jsonData, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			return config, errors.Wrap(err, fmt.Sprintf("Read config file error (%s)", configFilePath))
		}

		if err := json.Unmarshal(jsonData, &config); err != nil {
			return config, errors.Wrap(err, fmt.Sprintf("Unmarshal config file error (%s)", configFilePath))
		}
	}

	return
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	config, err := parseServerConfig()
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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router); err != nil {
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
