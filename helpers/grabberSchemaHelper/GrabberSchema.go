package grabberSchemaHelper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/d0kur0/webm-grabber/sources/fourChannel"
	"github.com/d0kur0/webm-grabber/sources/twoChannel"

	"github.com/d0kur0/webm-grabber/sources/types"
)

type grabberSchema struct {
	grabberSchema     []types.GrabberSchema
	allowedExtensions types.AllowedExtensions
	configFilePath    string
}

func (schema *grabberSchema) createConfigFromDefault() error {
	var grabberSchema = make(map[string][]types.Board, len(schema.grabberSchema))

	for _, schema := range schema.grabberSchema {
		grabberSchema[schema.Vendor.VendorName()] = schema.Boards
	}

	var configData = configStruct{
		AllowedExtensions: schema.allowedExtensions,
		GrabberSchema:     grabberSchema,
	}

	jsonBytes, err := json.MarshalIndent(configData, "", "\t")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Marshaling config error (%s)", schema.configFilePath))
	}

	if err := ioutil.WriteFile(schema.configFilePath, jsonBytes, 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Write in config data file error (%s)", schema.configFilePath))
	}

	return nil
}

func (schema *grabberSchema) parseConfig() (err error) {
	if _, err := os.Stat(schema.configFilePath); os.IsNotExist(err) {
		return schema.createConfigFromDefault()
	}

	jsonData, err := ioutil.ReadFile(schema.configFilePath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Read config file error (%s)", schema.configFilePath))
	}

	var configData configStruct
	if err := json.Unmarshal(jsonData, &configData); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Unmarshal config file error (%s)", schema.configFilePath))
	}

	schema.allowedExtensions = configData.AllowedExtensions
	for schemaIndex, schemaData := range schema.grabberSchema {
		if _, vendorsExists := configData.GrabberSchema[schemaData.Vendor.VendorName()]; !vendorsExists {
			schema.grabberSchema[schemaIndex] = schema.grabberSchema[len(schema.grabberSchema)-1]
			schema.grabberSchema[len(schema.grabberSchema)-1] = types.GrabberSchema{}
			schema.grabberSchema = schema.grabberSchema[:len(schema.grabberSchema)-1]

			continue
		}

		schema.grabberSchema[schemaIndex].Boards = configData.GrabberSchema[schemaData.Vendor.VendorName()]
	}

	return
}

func (schema *grabberSchema) Get() []types.GrabberSchema {
	return schema.grabberSchema
}

func Make() (schema *grabberSchema) {
	schema = &grabberSchema{configFilePath: "schema.json"}

	// Default settings
	schema.allowedExtensions = types.AllowedExtensions{".webm", ".mp4"}
	schema.grabberSchema = []types.GrabberSchema{
		{
			twoChannel.Make(schema.allowedExtensions),
			[]types.Board{
				{"b", "Бред"},
				{"vg", "Видео Игры"},
				{"a", "Аниме"},
				{"mu", "Музыка"},
				{"e", "Extreme Porn"},
				{"h", "Хентай"},
				{"fur", "Фурри"},
				{"kpop", "K-Pop"},
				{"asmr", "ASMR"},
			},
		},
		{
			fourChannel.Make(schema.allowedExtensions),
			[]types.Board{
				{"a", "Anime & Manga"},
				{"c", "Anime/Cute"},
				{"cgl", "Cosplay & EGL"},
				{"vg", "Video Game Generals"},
				{"co", "Comics & Cartoons"},
				{"g", "Technology"},
				{"b", "Random"},
				{"mu", "Music"},
				{"s", "Sexy Beautiful Women"},
				{"hc", "Hardcore"},
				{"h", "Hentai"},
				{"e", "Ecchi"},
				{"u", "Yuri"},
				{"d", "Hentai/Alternative"},
			},
		},
	}

	// Parse config file
	if err := schema.parseConfig(); err != nil {
		log.Println(err)
	}

	return
}
