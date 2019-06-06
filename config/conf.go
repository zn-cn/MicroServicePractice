package config

import (
	"io/ioutil"
	"log"
	"os"

	"fmt"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Env           string              `yaml:"Env"`
	ServiceName   map[string]string   `yaml:"ServiceName"`
	Version       map[string]string   `yaml:"Version"`
	DB            map[string]db       `yaml:"DB"`
	BrokerTopic   map[string]string   `yaml:"BrokerTopic"`
	BrokerAddrs   map[string][]string `yaml:"BrokerAddrs"`
	RegistryAddrs map[string][]string `yaml:"RegistryAddrs"`
	TracingAddr   map[string]string   `yaml:"TracingAddr"`
}

type db struct {
	DriverName  string `yaml:"DriverName"`
	Host        string `yaml:"Host"`
	Port        int32  `yaml:"Port"`
	DBName      string `yaml:"DBName"`
	User        string `yaml:"User"`
	PW          string `yaml:"PW"`
	AdminDBName string `yaml:"AdminDBName"`
}

var config Conf

func GetConfig() Conf {
	return config
}

func GetServiceName(service string) string {
	return config.ServiceName[service]
}

func GetVersion(service string) string {
	return config.Version[service]
}

func GetTracingAddr(service string) string {
	return config.TracingAddr[service]
}

func GetDB(service string) db {
	return config.DB[service]
}

func GetBrokerTopic(service string) string {
	return config.BrokerTopic[service]
}

func GetBrokerAddrs(service string) []string {
	if service == "" {
		service = "default"
	}
	return config.BrokerAddrs[service]
}

func GetRegistryAddrs(service string) []string {
	if service == "" {
		service = "default"
	}
	return config.RegistryAddrs[service]
}

func init() {
	env := os.Getenv("Env")
	if env == "" {
		env = "dev"
	}
	prefixPath := os.Getenv("PrefixPath")
	if prefixPath == "" {
		gopath := os.Getenv("GOPATH")
		prefixPath = gopath + "/src/Ethan/MicroServicePractice/"
	}
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s/config/conf-%s.yaml", prefixPath, env))
	if err != nil {
		log.Fatalf("read yaml config error: %v\n", err)
	}
	err = yaml.UnmarshalStrict(yamlFile, &config)
}
