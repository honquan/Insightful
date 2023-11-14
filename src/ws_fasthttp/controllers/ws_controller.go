package controllers

import (
	"context"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"insightful/src/ws_fasthttp/services"
	"log"
)

type WebsocketController interface {
	WebsocketWorkerFastHttp(ctx *fasthttp.RequestCtx)
}

type websocketController struct {
	websocketService services.WebsocketService
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  10240000,
	WriteBufferSize: 1024,
}

func NewWebsocketController() WebsocketController {
	var websocketService services.WebsocketService
	_ = services.GetServiceContainer().Invoke(func(s services.WebsocketService) {
		websocketService = s
	})

	return &websocketController{
		websocketService: websocketService,
	}
}

func (s *websocketController) WebsocketWorkerFastHttp(ctx *fasthttp.RequestCtx) {
	upgrader.CheckOrigin = func(ctx *fasthttp.RequestCtx) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		// listen indefinitely for new messages coming
		// through on our WebSocket connection
		_ = s.websocketService.ReaderFastHttpWithAnts(context.Background(), conn)
	})
	if err != nil {
		log.Println(err)
		return
	}

}
