package main

import (
	"insightful/src/apis/pkg/worker"
	"insightful/src/apis/router"
	"insightful/src/apis/services"
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
	services.InitialServices()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// run custom job worker without redis
	//custom_worker.JobQueue = make(chan custom_worker.Job, conf.EnvConfig.MaxWorker)
	//dispatcher := custom_worker.NewDispatcher(conf.EnvConfig.MaxWorker)
	//dispatcher.Run()

	// run muster
	//muster.Run()

	// init router
	a := router.App{}
	a.InitRouter()

	// run worker go worker
	go worker.RunGoWorker()

	// run go craft
	go worker.RunGoCraft()

	// run
	a.Run(":8899")
}
