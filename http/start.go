package http

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func Start() error {
	http.HandleFunc("/schema", getSchema)
	http.HandleFunc("/files", getFiles)
	http.HandleFunc("/filesByCondition", getFilesByCondition)

	return http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), nil)
}
