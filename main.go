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
	log.Println("Update tick done.")
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

	schemaAsBytes, err := json.Marshal(schema)
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
	fmt.Println(`
******************************************** WEBM API ********************************************
*                                                                                                *
*  HTTP Endpoints:                                                                               *
*                                                                                                *
*    Get current server grabber schema: [request type: GET]                                      *
*    http://localhost:3000/schema                                                                *
*                                                                                                *
*    Get all grabbed files: [request type: GET]                                                  *
*    http://localhost:3000/files                                                                 *
*                                                                                                *
*    Get grabbed files by specific vendor and boards: [request type: POST]                       *
*    http://localhost:3000/filesWithCondition                                                    *
*        This request need body with condition struct:                                           *
*                { "<vendor name from grabber schema>": ["<board name1>", "<board name 2>"] }    *
*                For example:                                                                    *
*                { "2ch": ["b", "media", "vg"], "4chan": ["b"] }                                 *
*                                                                                                *
**************************************************************************************************
	`)

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}

	twoCh := twoChannel.Make(types.AllowedExtensions{})
	fourCh := fourChannel.Make(types.AllowedExtensions{})

	for vendor, boards := range config.Schema {
		var filledBoards []types.Board
		for _, board := range boards {
			filledBoards = append(filledBoards, types.Board{
				Name:        board,
				Description: "any",
			})
		}

		var vendorInstance *types.VendorInterface
		if vendor == TwoChName {
			vendorInstance = &twoCh
		}
		if vendor == FourChName {
			vendorInstance = &fourCh
		}

		if vendorInstance == nil {
			log.Fatalln("Undefined vendor from config file")
			return
		}

		schema = append(schema, types.GrabberSchema{Vendor: *vendorInstance, Boards: filledBoards})
	}

	log.Println("> First grabbing files, before startup http server...")
	output = webmGrabber.GrabberProcess(schema)
	log.Println("> Files grabbed.")
	schedule(updateTick, 1*time.Minute)

	http.HandleFunc("/schema", HttpGetSchema)
	http.HandleFunc("/files", HttpGetFiles)

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
