package services

import (
	"go.uber.org/dig"
)

// serviceContainer is a global ServiceProvider.
var serviceContainer *dig.Container

func InitServices() {
	container := dig.New()

	serviceContainer = container
}

// GetServiceContainer return a new instance of ServiceContainer
func GetServiceContainer() *dig.Container {
	return serviceContainer
}
