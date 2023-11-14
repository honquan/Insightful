package services

import (
	"go.uber.org/dig"
	"insightful/common/config"
	"insightful/connection"
	repository "insightful/src/apis/repositories"
)

// serviceContainer is a global ServiceProvider.
var serviceContainer *dig.Container

func InitialServices() {
	container := dig.New()
	_ = container.Provide(config.NewConfig)

	// provide connect mongo
	_ = container.Provide(connection.InitMongo)

	// provide repo
	_ = container.Provide(repository.NewInsightfullRepository)

	// provide service
	_ = container.Provide(NewWebsocketService)

	serviceContainer = container
}

// GetServiceContainer return a new instance of ServiceContainer
func GetServiceContainer() *dig.Container {
	return serviceContainer
}
