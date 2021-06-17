package util

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"strconv"
)


var configFilePath = "config.ini"

type ConfigINI struct {
	ServerHost 		string
	ServerPort		int
}

func LoadConfig() (configINI *ConfigINI, err error) {
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Println(fmt.Sprintf("Fail to read file: %v", err))
		return nil, err
	}

	serverPortNumber, err := strconv.Atoi(cfg.Section("server").Key("port").String())

	configINI = &ConfigINI{
		ServerHost: cfg.Section("server").Key("host").String(),
		ServerPort: serverPortNumber,
	}

	return configINI, err
}