package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ApiGatewayConfig struct {
	Version string `json:"version"`

	Debug       int    `json:"debug"`
	FileVersion int    `json:"file_version"`
	Banner      string `json:"banner"`

	ListenAddr string `json:"Listen_addr"`
	ListenPort int    `json:"Listen_port"`

	PluginList []string `json:"plugin_list"`

	ESUrl string `json:"es_url"`

	ApimapConfig string `json:"apimap_config"`
}

func ParseConfig(filepath string, config *ApiGatewayConfig) error {

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, config)
	if err != nil {
		return err
	}
	return nil
}
