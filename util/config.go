package util

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"path"
	"strconv"
)


const configFileName = "config.ini"

type ConfigINI struct {
	ServerHost 		string
	ServerPort		int
	Protocol		string
}

func LoadConfig(homePath string) (configINI *ConfigINI, err error) {
	cfg, err := ini.Load(path.Join(homePath, configFileName))
	if err != nil {
		log.Println(fmt.Sprintf("Fail to read file: %v", err))
		return nil, err
	}

	serverPortNumber, err := strconv.Atoi(cfg.Section("server").Key("port").String())

	configINI = &ConfigINI{
		ServerHost: cfg.Section("server").Key("host").String(),
		ServerPort: serverPortNumber,
		Protocol: cfg.Section("common").Key("protocol").String(),
	}

	return configINI, err
}