package service

import (
	"go.uber.org/dig"
	"insightful/src/apis/conf"
	worker "insightful/src/apis/kit/custom_worker"
)

// serviceContainer is a global ServiceProvider.
var serviceContainer *dig.Container

func InitialServices() {
	container := dig.New()

	_ = container.Provide(func() *worker.Dispatcher {
		return worker.NewDispatcher(conf.EnvConfig.MaxWorker).Run()
	})

	_ = container.Provide(NewWebsocketService)

	serviceContainer = container
}

// GetServiceContainer return a new instance of ServiceContainer
func GetServiceContainer() *dig.Container {
	return serviceContainer
}
