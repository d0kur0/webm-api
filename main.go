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
	Port           int
	Schema         map[string][]string
	Extensions     types.AllowedExtensions
	UpdateInterval int
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
}

type SimplifySchemaItem struct {
	Vendor string   `json:"vendor"`
	Boards []string `json:"board"`
}

func HttpGetSchema(w http.ResponseWriter, r *http.Request) {
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

type GetFilesCondition map[string][]string

func HttpGetFilesByCondition(w http.ResponseWriter, r *http.Request) {
	var condition GetFilesCondition
	err := json.NewDecoder(r.Body).Decode(&condition)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var filteredFiles types.Output
	for _, outItem := range output {
		boards, isDesiredVendor := condition[outItem.VendorName]
		if !isDesiredVendor {
			continue
		}

		isDesiredBoard := false
		for _, board := range boards {
			if board == outItem.BoardName {
				isDesiredBoard = true
				break
			}
		}

		if !isDesiredBoard {
			continue
		}

		filteredFiles = append(filteredFiles, outItem)
	}

	outputAsBytes, err := json.Marshal(filteredFiles)
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

		var vendorInstances = make(map[string]types.VendorInterface)
		vendorInstances[TwoChName] = twoChannel.Make(config.Extensions)
		vendorInstances[FourChName] = fourChannel.Make(config.Extensions)

		vendorInstance, isExists := vendorInstances[vendor]
		if !isExists {
			log.Fatalln("Undefined vendor from config file")
			return
		}

		schema = append(schema, types.GrabberSchema{Vendor: vendorInstance, Boards: filledBoards})
	}

	schedule(updateTick, time.Duration(config.UpdateInterval)*time.Minute)
	output = webmGrabber.GrabberProcess(schema)

	http.HandleFunc("/schema", HttpGetSchema)
	http.HandleFunc("/files", HttpGetFiles)
	http.HandleFunc("/filesByCondition", HttpGetFilesByCondition)

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
