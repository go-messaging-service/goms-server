package technicalServices

import (
	"encoding/json"
	"goMS/src/technical/material"
	"goMS/src/technical/services/logger"
	"io/ioutil"
)

// simplify the way of access to structs
//type Config technicalMaterial.Config
//type TopicConfig technicalMaterial.TopicConfig
//type ServerConfig technicalMaterial.ServerConfig

type ConfigLoader struct {
	topicConfig  *technicalMaterial.TopicConfig
	serverConfig *technicalMaterial.ServerConfig
}

// LoadTopics reads the config file for the topics and fills the TopicConfig field.
// The default location is /config/topics.json
func (cl *ConfigLoader) loadTopics(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		logger.Error("Error reading " + filename)
		logger.Fatal(err.Error())
	}

	cl.topicConfig = &technicalMaterial.TopicConfig{}
	json.Unmarshal(data, cl.topicConfig)
}

// LoadConfig loads the given server config file.
// The default location is /config/server.json
func (cl *ConfigLoader) LoadConfig(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		logger.Error("Error reading " + filename)
		logger.Fatal(err.Error())
	}

	cl.serverConfig = &technicalMaterial.ServerConfig{}
	json.Unmarshal(data, cl.serverConfig)

	cl.loadTopics(cl.serverConfig.TopicLocation)
}

// GetConfig creates a configuration every topic and server information is in.
func (cl *ConfigLoader) GetConfig() technicalMaterial.Config {
	return technicalMaterial.Config{
		TopicConfig:  *cl.topicConfig,
		ServerConfig: *cl.serverConfig,
	}
}
