package settingsManager

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func Parse() (config ServerConfig, err error) {

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
