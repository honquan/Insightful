package main

import (
	"insightful/common/logger"
	"insightful/src/apis/conf"
	"insightful/src/apis/router"
	"log"
	"strings"
)

func init() {
	// Init logging
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      strings.ToLower(conf.EnvConfig.LogLevel),
		ConsoleJSONFormat: true,
		EnableFile:        false,
	}
	logger := logger.NewLogger(config, logger.InstanceZapLogger)
	if logger == nil {
		log.Printf("Could not instantiate log")
	}
}

func main() {
	a := router.App{}
	a.InitRouter()

	a.Run(":8888")
}
