package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(InitCmd)
}

type CfgStruct struct {
	Port           int                 `json:"port"`
	UpdateInterval int                 `json:"updateInterval"`
	Schema         map[string][]string `json:"schema"`
	Extensions     []string            `json:"extensions"`
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init base configuration",
	Run: func(cmd *cobra.Command, args []string) {
		emptyCfg := CfgStruct{}

		jsonBytes, err := json.MarshalIndent(emptyCfg, "", " ")
		if err != nil {
			log.Fatalln(err)
		}

		err = os.WriteFile(".webm-api.json", jsonBytes, 0644)
		if err != nil {
			log.Fatalln(err)
		}
	},
}
