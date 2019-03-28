package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hauke96/sigolo"
)

// simplify the way of access to structs
//type Config technicalMaterial.Config
//type TopicConfig technicalMaterial.TopicConfig
//type ServerConfig technicalMaterial.ServerConfig

type ConfigLoader struct {
	topicConfig  *TopicConfig
	serverConfig *ServerConfig
}

// LoadTopics reads the config file for the topics and fills the TopicConfig field.
// The default location is /config/topics.json
func (cl *ConfigLoader) loadTopics(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		sigolo.Error("Error reading " + filename)
		sigolo.Error("\n\nAhhh, *urg*, I'm sorry but there was a really bad error inside of me. Above the stack trace is a message marked with [FATAL], you'll find some information there.\nIf not, feel free to contact my maker via:\n\n    goms@hauke-stieler.de\n\nI hope my death ..    . eh ... crash is only an exception and will be fixed soon ... my power ... leaves me ... good bye ... x.x")
		sigolo.Error(err.Error())
	}

	cl.topicConfig = &TopicConfig{}
	json.Unmarshal(data, cl.topicConfig)
}

// LoadConfig loads the given server config file.
// The default location is /config/server.json
func (cl *ConfigLoader) LoadConfig(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		sigolo.Error("Error reading " + filename)
		sigolo.Error("\n\nAhhh, *urg*, I'm sorry but there was a really bad error inside of me. Above the stack trace is a message marked with [FATAL], you'll find some information there.\nIf not, feel free to contact my maker via:\n\n    goms@hauke-stieler.de\n\nI hope my death ..    . eh ... crash is only an exception and will be fixed soon ... my power ... leaves me ... good bye ... x.x")
		sigolo.Fatal(err.Error())
	}

	cl.serverConfig = &ServerConfig{}
	json.Unmarshal(data, cl.serverConfig)

	cl.loadTopics(cl.serverConfig.TopicLocation)
}

// GetConfig creates a configuration every topic and server information is in.
func (cl *ConfigLoader) GetConfig() Config {
	return Config{
		TopicConfig:  *cl.topicConfig,
		ServerConfig: *cl.serverConfig,
	}
}
