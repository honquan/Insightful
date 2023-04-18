package main

import (
	"insightful/src/apis/pkg/worker"
	"insightful/src/apis/router"
	"log"
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
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	a := router.App{}
	// init router
	a.InitRouter()
	// run worker go worker
	go worker.RunGoWorker()

	// run go craft
	go worker.RunGoCraft()

	// run
	a.Run("127.0.0.1:8899")
}
