package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/d0kur0/webm-api/http"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(StartCmd)
}

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Server starting at: http://localhost:%d", viper.GetInt("port")))
		err := http.Start()
		if err != nil {
			fmt.Println(errors.Wrap(err, "failed starting http server"))
		}
	},
}
