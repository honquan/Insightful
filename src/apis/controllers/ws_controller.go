package controllers

import (
	"github.com/gocraft/work"
	"github.com/gorilla/websocket"
	go_worker "insightful/src/apis/pkg/worker"
	"log"
	"net/http"
	"time"
)

type WsController struct {
	BaseController
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *WsController) WebsocketWorker(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	// helpful log statement to show connections
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	readerWithGoCraft(ws)
	//readerWithGoWorker(ws)
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func readerWithGoCraft(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println("Client said: ", string(p))
		// enqueue go craft
		enqueueJobCraft(
			"saveCoordinateCassandra",
			work.Q{"message": string(p)},
		)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func readerWithGoWorker(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println("Client said: ", string(p))
		go_worker.AddJob("Sample", time.Now().UTC(), string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
