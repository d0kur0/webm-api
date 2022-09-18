package util

import (
	"log"

	"github.com/spf13/viper"

	"github.com/d0kur0/webm-grabber/types"
	"github.com/d0kur0/webm-grabber/vendors/fourChannel"
	"github.com/d0kur0/webm-grabber/vendors/twoChannel"
)

func ParseSchema() (schema []types.GrabberSchema) {
	for vendor, boards := range viper.GetStringMapStringSlice("schema") {
		var filledBoards []types.Board
		for _, board := range boards {
			filledBoards = append(filledBoards, types.Board{
				Name:        board,
				Description: "any",
			})
		}

		var vendorInstances = make(map[string]types.VendorInterface)
		vendorInstances[twoChannel.VendorName] = twoChannel.Make(viper.GetStringSlice("extensions"))
		vendorInstances[fourChannel.VendorName] = fourChannel.Make(viper.GetStringSlice("extensions"))

		vendorInstance, isExists := vendorInstances[vendor]
		if !isExists {
			log.Fatalln("Undefined vendor from config file")
			return
		}

		schema = append(schema, types.GrabberSchema{Vendor: vendorInstance, Boards: filledBoards})
	}

	return schema
}
