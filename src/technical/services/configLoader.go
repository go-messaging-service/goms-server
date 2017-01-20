package technicalServices

import (
	"encoding/json"
	"goMS/src/technical/material"
	"goMS/src/technical/services/logger"
	"io/ioutil"
)

// simplify the way of access to structs
type Config technicalMaterial.Config
type TopicConfig technicalMaterial.TopicConfig
type ServerConfig technicalMaterial.ServerConfig

type ConfigLoader struct {
	TopicConfig  *TopicConfig
	ServerConfig *ServerConfig
}

// LoadTopics reads the config file for the topics and fills the TopicConfig field
// The default location is /config/topics.json
func (cl *ConfigLoader) loadTopics(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		logger.Error("Error reading " + filename)
		logger.Fatal(err.Error())
	}

	cl.TopicConfig = &TopicConfig{}
	json.Unmarshal(data, cl.TopicConfig)
}

func (cl *ConfigLoader) LoadConfig(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		logger.Error("Error reading " + filename)
		logger.Fatal(err.Error())
	}

	cl.ServerConfig = &ServerConfig{}
	json.Unmarshal(data, cl.ServerConfig)

	cl.loadTopics(cl.ServerConfig.TopicLocation)
}
