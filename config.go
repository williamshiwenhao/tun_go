package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/withmandala/go-log"
)

type config struct {
	ToAddr   string `json:"to_addr"`
	SelfPort int    `json:"self_port"`
}

// Config global config
var Config *config
var logger = log.New(os.Stderr).WithColor()

const configPath = "config.json"

func init() {
	Config = &config{}
	readConfig(configPath, Config)
}

// InitConfig init the global config
func readConfig(path string, config interface{}) {
	configFd, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer configFd.Close()
	rawData, err := ioutil.ReadAll(configFd)
	if err != nil {
		logger.Fatalf("Cannot read from config file, err: %+v", err)
	}
	err = json.Unmarshal(rawData, config)
	if err != nil {
		logger.Fatalf("Cannot parse init config, err: %+v", err)
	}
	logOutput, _ := json.MarshalIndent(config, "", "\t")
	logger.Infof("[Init] Config parsed! [Config = \n%v", string(logOutput))
}
