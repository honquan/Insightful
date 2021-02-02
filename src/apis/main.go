package main

import (
	"insightful/src/apis/router"
)

func init() {
	// Init logging
	//config := logger.Configuration{
	//	EnableConsole:     true,
	//	ConsoleLevel:      strings.ToLower(conf.EnvConfig.LogLevel),
	//	ConsoleJSONFormat: true,
	//	EnableFile:        false,
	//}
	//logger := logger.NewLogger(config, logger.InstanceZapLogger)
	//if logger == nil {
	//	log.Printf("Could not instantiate log")
	//}

	// Init services
	//services.InitServices()
}

func main() {
	a := router.App{}
	a.InitRouter()

	//a.Run(":8888")
	a.Run("127.0.0.1:8889")
}
