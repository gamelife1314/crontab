package crond

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// HttpConfig set http service configuration
type HttpConfig struct {
	Address      string `yaml:"address" json:"address"`
	Port         int    `yaml:"port" json:"port"`
	ReadTimeout  int    `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout" json:"writeTimeout"`
}

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

// ServiceConfig set crond service configuration
type ServiceConfig struct {
	Http  HttpConfig  `yaml:"http" json:"http"`
	Etcd  EtcdConfig  `yaml:"etcd" json:"etcd"`
	Mongo MongoConfig `yaml:"mongo" json:"mongo"`
}

var Config *ServiceConfig

// LoadConfig load service configuration from yaml file
func LoadConfig(filename string) (err error) {
	var (
		content []byte
		config  ServiceConfig
	)
	if filename == "" {
		filename = os.Getenv("CROND_CONFIG_FILE")
	}
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	config = ServiceConfig{}
	if err = yaml.Unmarshal(content, &config); err != nil {
		return
	}
	Config = &config
	return
}

func (s *ServiceConfig) String() string {
	if content, err := json.Marshal(s); err != nil {
		return ""
	} else {
		return string(content)
	}
}
