package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	// define application mode as an environment variable
	ENV_MODE = "ENVIRONMENT_MODE"

	DevelopmentMode = "dev"
	ProductionMode  = "prod"
)

type HttpConf struct {
	BindIp      string `yaml:"bindIp"`
	BindPort    int    `yaml:"bindPort"`
	ContextPath string `yaml:"contextPath"`
}

type ServerConf struct {
	Http *HttpConf `yaml:"http"`
}

type AppConf struct {
	Server *ServerConf `yaml:"server"`
}

var Server *ServerConf

// application mode either 'release' or 'dev'
var EnvironmentMode string = DevelopmentMode

var defaultConfig = &AppConf{
	Server: &ServerConf{
		Http: &HttpConf{
			BindIp:      "127.0.0.1",
			BindPort:    80,
			ContextPath: "",
		},
	},
}

func init() {
	EnvironmentMode = os.Getenv(ENV_MODE)
	if len(EnvironmentMode) == 0 {
		EnvironmentMode = DevelopmentMode
	}
}

func Initialize() {

	ymlPath := flag.String("config", "config.yml", "yaml based configuration path")
	flag.Parse()

	var appConf AppConf

	// check if file exist
	if _, err := os.Stat(*ymlPath); errors.Is(err, os.ErrNotExist) {
		log.Println("config file don't exist. default config is applied")
		Server = defaultConfig.Server
		return
	}

	ymlConfig, err := os.ReadFile(*ymlPath)
	if err != nil {
		log.Fatalf("Could not read config file. error: %v", err)
	}

	err = yaml.Unmarshal(ymlConfig, &appConf)
	if err != nil {
		log.Fatalf("Could not unmarshal config file. error: %v", err)
	}

	prepareFinalConfig(&appConf)
	Server = appConf.Server
}

func prepareFinalConfig(appConfig *AppConf) {

	if appConfig.Server == nil {
		appConfig.Server = defaultConfig.Server
		return
	}

	if len(appConfig.Server.Http.BindIp) == 0 {
		appConfig.Server.Http.BindIp = defaultConfig.Server.Http.BindIp
	}

	if appConfig.Server.Http.BindPort == 0 {
		appConfig.Server.Http.BindPort = defaultConfig.Server.Http.BindPort
	}

	if len(appConfig.Server.Http.ContextPath) == 0 {
		appConfig.Server.Http.ContextPath = defaultConfig.Server.Http.ContextPath
	}
}
