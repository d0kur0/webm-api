package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	webmGrabber "github.com/d0kur0/webm-grabber"

	"github.com/d0kur0/webm-grabber/vendors/fourChannel"

	"github.com/d0kur0/webm-grabber/types"
	"github.com/d0kur0/webm-grabber/vendors/twoChannel"
	"github.com/pkg/errors"
)

const ConfigPath = "./config.json"
const TwoChName = "2ch"
const FourChName = "4chan"

var schema []types.GrabberSchema
var output types.Output

type Config struct {
	Port       int
	Schema     map[string][]string
	Extensions types.AllowedExtensions
}

func readConfig() (config Config, err error) {
	jsonData, err := os.ReadFile(ConfigPath)
	if err != nil {
		return config, errors.Wrap(err, fmt.Sprintf("Read config file error (%s)", ConfigPath))
	}

	if err := json.Unmarshal(jsonData, &config); err != nil {
		return config, errors.Wrap(err, fmt.Sprintf("Unmarshal config file error (%s)", ConfigPath))
	}

	return
}

func schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func updateTick() {
	output = webmGrabber.GrabberProcess(schema)
	log.Println("Update tick done.", output)
}

func HttpGetSchema(w http.ResponseWriter, r *http.Request) {
	type SimplifySchemaItem struct {
		Vendor string
		Boards []string
	}

	var simplifySchema []SimplifySchemaItem
	for _, schemaEl := range schema {
		var boards []string
		for _, board := range schemaEl.Boards {
			boards = append(boards, board.Name)
		}
		simplifySchema = append(simplifySchema, SimplifySchemaItem{Vendor: schemaEl.Vendor.VendorName(), Boards: boards})
	}

	schemaAsBytes, err := json.Marshal(simplifySchema)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(w, string(schemaAsBytes))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func HttpGetFiles(w http.ResponseWriter, r *http.Request) {
	outputAsBytes, err := json.Marshal(output)
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

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}

	for vendor, boards := range config.Schema {
		var filledBoards []types.Board
		for _, board := range boards {
			filledBoards = append(filledBoards, types.Board{
				Name:        board,
				Description: "any",
			})
		}

		var vendorInstances map[string]types.VendorInterface
		vendorInstances[TwoChName] = twoChannel.Make(types.AllowedExtensions{})
		vendorInstances[FourChName] = fourChannel.Make(types.AllowedExtensions{})

		vendorInstance, isExists := vendorInstances[vendor]
		if !isExists {
			log.Fatalln("Undefined vendor from config file")
			return
		}

		schema = append(schema, types.GrabberSchema{Vendor: vendorInstance, Boards: filledBoards})
	}

	schedule(updateTick, 5*time.Minute)

	http.HandleFunc("/schema", HttpGetSchema)
	http.HandleFunc("/files", HttpGetFiles)

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
