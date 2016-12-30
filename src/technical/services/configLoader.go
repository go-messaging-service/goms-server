package services

import (
	"encoding/json"
	"goMS/src/technical/material"
	"goMS/src/technical/services/logger"
	"io/ioutil"
	"os"
)

// simplify the way of access to structs
type Config material.Config
type TopicConfig material.TopicConfig

type ConfiLoader struct {
	topicConfig TopicConfig
}

func (cl *ConfiLoader) LoadTopics(filename string) {
	cl.topicConfig = TopicConfig{}
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(data, cl.topicConfig)
}
