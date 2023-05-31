package controllers

import (
	"context"
	"github.com/gorilla/websocket"
	"insightful/src/apis/services"
	"log"
	"net/http"
)

type WebsocketController struct {
	BaseController
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *WebsocketController) WebsocketWorkerGoCraft(w http.ResponseWriter, r *http.Request) {
	var websocketService services.WebsocketService
	_ = services.GetServiceContainer().Invoke(func(s services.WebsocketService) {
		websocketService = s
	})

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	// helpful log statement to show connections
	//err = ws.WriteMessage(1, []byte("Hi Client!"))
	//if err != nil {
	//	log.Println(err)
	//}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	err = websocketService.ReaderWithGoCraft(context.Background(), ws)
}

func (s *WebsocketController) WebsocketWorkerGoWorker(w http.ResponseWriter, r *http.Request) {
	var websocketService services.WebsocketService
	_ = services.GetServiceContainer().Invoke(func(s services.WebsocketService) {
		websocketService = s
	})

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	// helpful log statement to show connections
	//err = ws.WriteMessage(1, []byte("Hi Client!"))
	//if err != nil {
	//	log.Println(err)
	//}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	err = websocketService.ReaderWithGoWorker(context.Background(), ws)
}

func (s *WebsocketController) WebsocketWorkerPool(w http.ResponseWriter, r *http.Request) {
	var websocketService services.WebsocketService
	_ = services.GetServiceContainer().Invoke(func(s services.WebsocketService) {
		websocketService = s
	})

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	err = websocketService.ReaderWithCustomWorkerPool(context.Background(), ws)
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
//func (s *WebsocketController) readerWithGoCraft(conn *websocket.Conn) {
//	for {
//		// read in a message
//		messageType, p, err := conn.ReadMessage()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//
//		// enqueue go craft
//		enqueueJobCraft(
//			enum.JobNameCoordinate,
//			work.Q{enum.GoCraftMessage: p},
//		)
//
//		if err := conn.WriteMessage(messageType, p); err != nil {
//			log.Println(err)
//			return
//		}
//
//	}
//}

//func (s *WebsocketController) readerWithGoWorker(conn *websocket.Conn) {
//	for {
//		// read in a message
//		messageType, p, err := conn.ReadMessage()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//
//		go_worker.AddJob(enum.JobNameCoordinate, time.Now().UTC(), p)
//
//		if err := conn.WriteMessage(messageType, p); err != nil {
//			log.Println(err)
//			return
//		}
//
//	}
//}
