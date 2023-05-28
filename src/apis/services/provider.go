package services

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"insightful/src/apis/conf"
	worker "insightful/src/apis/kit/custom_worker"
	"insightful/src/apis/pkg/connection"
	repository "insightful/src/apis/repositories"
)

// serviceContainer is a global ServiceProvider.
var serviceContainer *dig.Container

func InitialServices() {
	container := dig.New()

	_ = container.Provide(func() *worker.Dispatcher {
		return worker.NewDispatcher(conf.EnvConfig.MaxWorker).Run()
	})

	// provide connect mongo
	_ = container.Provide(func() *mongo.Database {
		mongoConnection, _, err := connection.NewMongoConnection()
		if err != nil {
			panic(err)
		}
		//defer closeFunc()

		return mongoConnection
	})

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
