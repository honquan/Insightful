package main

import (
	"insightful/src/worker/cmd"
	"insightful/src/worker/services"
)

func init() {
	// init services
	services.InitServices()
}

func main() {
	cmd.Execute()
}
