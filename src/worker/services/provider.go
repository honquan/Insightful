package services

import (
	"go.uber.org/dig"
	"insightful/common/config"
	"insightful/connection"
	"insightful/src/worker/repositories"
)

// serviceContainer is a global ServiceProvider.
var serviceContainer *dig.Container

func InitServices() {
	container := dig.New()

	_ = container.Provide(config.NewConfig)
	_ = container.Provide(connection.InitPostgres)
	_ = container.Provide(connection.InitMongo)

	_ = container.Provide(repositories.NewPostgresInsightfullRepository)
	_ = container.Provide(repositories.NewMongoInsightfullRepository)

	_ = container.Provide(NewAnalyzeDataService)

	serviceContainer = container
}

// GetServiceContainer return a new instance of ServiceContainer
func GetServiceContainer() *dig.Container {
	return serviceContainer
}
