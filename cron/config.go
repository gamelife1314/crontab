package cron

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// EtcdConfig set etcd connect configuration
type EtcdConfig struct {
	Endpoints   []string `yaml:"endPoints" json:"endpoints"`
	DialTimeout int      `yaml:"dialTimeout" json:"dialTimeout"`
}

// MongoConfig set mongo connect configuration
type MongoConfig struct {
	Url        string `yaml:"url" json:"url"`
	Database   string `yaml:"database" json:"database"`
	Collection string `yaml:"collection" json:"collection"`
}

type ClientConfig struct {
	JobLogBatchSize     int    `json:"jobLogBatchSize" yaml:"jobLogBatchSize"`
	JobLogCommitTimeout int    `json:"jobLogCommitTimeout" yaml:"jobLogCommitTimeout"`
	BashExecutePath     string `json:"bashExecutePath" yaml:"bashExecutePath"`
}

type ServiceConfig struct {
	Etcd   EtcdConfig   `json:"etcd" yaml:"etcd"`
	Mongo  MongoConfig  `yaml:"mongo" json:"mongo"`
	Client ClientConfig `yaml:"client" json:"client"`
}

var GlobalConfig *ServiceConfig

func InitConfig(filename string) (err error) {
	var (
		content []byte
		config  ServiceConfig
	)

	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	config = ServiceConfig{}
	if err = yaml.Unmarshal(content, &config); err != nil {
		return
	}

	GlobalConfig = &config

	return nil
}

func (s *ServiceConfig) String() string {
	if content, err := json.Marshal(s); err != nil {
		return ""
	} else {
		return string(content)
	}
}
