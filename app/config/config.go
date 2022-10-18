package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	developmentEnv = "dev"
	productionEnv  = "prod"
)

// either 'prod' or 'dev'
var environment string

type HttpConf struct {
	BindIp      string `yaml:"bindIp"`
	BindPort    int    `yaml:"bindPort"`
	ContextPath string `yaml:"contextPath"`
}

type ServerConf struct {
	Http *HttpConf `yaml:"http"`
}

type AppConf struct {
	Server    *ServerConf    `yaml:"server"`
	Datastore *DatastoreConf `yaml:"datastore"`
}

type DatastoreConf struct {
	Mongo *MongoConf `yaml:"mongo"`
}

type MongoConf struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	App *AppConf `yaml:"App"`
}

var Server *ServerConf
var Datastore *DatastoreConf

func init() {
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case developmentEnv:
		environment = developmentEnv
		return
	case productionEnv:
		environment = productionEnv
		return
	}
	environment = developmentEnv
}

func IsProduction() bool {
	return environment == productionEnv
}

func Initialize() {
	ymlPath := flag.String("config", "config.yml", "yaml based configuration path")
	flag.Parse()

	var conf Config

	// check if file exist
	if _, err := os.Stat(*ymlPath); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("config file don't exist.")
	}

	ymlConfig, err := os.ReadFile(*ymlPath)
	if err != nil {
		log.Fatalf("Could not read config file. error: %v\n", err)
	}

	err = yaml.Unmarshal(ymlConfig, &conf)
	if err != nil {
		log.Fatalf("Could not unmarshal config file. error: %v\n", err)
	}

	Server = conf.App.Server
	Datastore = conf.App.Datastore
}
