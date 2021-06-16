package util

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
)


var configFilePath string = "config.ini"

func LoadConfig() (serverURL string) {
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to read file: %v", err))
	}

	serverHost := cfg.Section("server").Key("host").String()
	serverPort := cfg.Section("server").Key("port").String()

	serverURL = fmt.Sprintf(fmt.Sprintf("%v:%v", serverHost, serverPort))

	return serverURL
}