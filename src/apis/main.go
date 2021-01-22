package main

import (
	"insightful/src/apis/router"
)

func main() {
	a := router.App{}
	a.InitRouter()

	a.Run("127.0.0.1:8899")
}
