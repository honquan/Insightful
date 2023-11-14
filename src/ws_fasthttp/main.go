package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"insightful/src/ws_fasthttp/controllers"
	"insightful/src/ws_fasthttp/services"
	"log"
	_ "net/http/pprof"
)

func init() {
	// Init services
	services.InitialServices()
}

var addr = flag.String("addr", "localhost:8899", "http service address")

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	wsController := controllers.NewWebsocketController()
	//_ = fasthttp.ListenAndServe(":8899/apis/fasthttp-ws/worker-ants", wsController.WebsocketWorkerFastHttp)

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/apis/fasthttp-ws/worker-ants":
			wsController.WebsocketWorkerFastHttp(ctx)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}

	server := fasthttp.Server{
		Name:    "EchoExample",
		Handler: requestHandler,
	}

	log.Fatal(server.ListenAndServe(*addr))
}
