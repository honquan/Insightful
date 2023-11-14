package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"insightful/src/apis/controllers"
	"insightful/src/apis/dtos"
	"log"
	"net/http"
)

const (
	RouterWSPrefix = "/apis"
)

type App struct {
	Router  *mux.Router
	MysqlDB *gorm.DB
}

func (a *App) InitRouter() {
	// creates a new instance of a mux router
	a.Router = mux.NewRouter().StrictSlash(true)

	// add router health check controller
	healthController := &controllers.HealthController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/health", RouterWSPrefix), healthController.HealthCheck).Methods(http.MethodGet)
	a.Router.HandleFunc(fmt.Sprintf("%v/build", RouterWSPrefix), healthController.BuildData).Methods(http.MethodPost)

	// custom job worker
	customJobController := &controllers.CustomJobController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/custom/worker", RouterWSPrefix), customJobController.WorkerJobCustom).Methods(http.MethodPost)

	// add router websocket controller
	wsController := controllers.NewWebsocketController()
	a.Router.HandleFunc(fmt.Sprintf("%v/ws/worker-craft", RouterWSPrefix), wsController.WebsocketWorkerGoCraft).Methods(http.MethodGet)
	a.Router.HandleFunc(fmt.Sprintf("%v/ws/worker-go", RouterWSPrefix), wsController.WebsocketWorkerGoWorker).Methods(http.MethodGet)
	a.Router.HandleFunc(fmt.Sprintf("%v/ws/worker-pool", RouterWSPrefix), wsController.WebsocketWorkerPool).Methods(http.MethodGet)
	a.Router.HandleFunc(fmt.Sprintf("%v/ws/worker-ants", RouterWSPrefix), wsController.WebsocketAntsWorker).Methods(http.MethodGet)

	// normal job
	normalController := &controllers.NormalJobController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/job-normal/go-worker", RouterWSPrefix), normalController.NormalJobWorkerGoWorker).Methods(http.MethodPost)
	a.Router.HandleFunc(fmt.Sprintf("%v/job-normal/craft-worker", RouterWSPrefix), normalController.NormalJobWorkerGoCraft).Methods(http.MethodPost)

	// job
	jobController := &controllers.JobController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/job/worker-goworker", RouterWSPrefix), jobController.WorkerGoWorker).Methods(http.MethodGet)
	a.Router.HandleFunc(fmt.Sprintf("%v/job/worker-gocraft", RouterWSPrefix), jobController.WorkerGoCraft).Methods(http.MethodGet)

	// register middleware
	//a.Router.Use(a.recoverPanicMiddleware)
}

func (a *App) Run(addr string) {
	log.Printf("Starting server at %v", addr)
	err := http.ListenAndServe(fmt.Sprintf("%v", addr), a.Router)
	if err != nil {
		panic(err)
	}
}

func (a *App) recoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Printf("Panic handling route \"%v\", details: %v", r.URL.EscapedPath(), err)

				a.serveResponse(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *App) serveResponse(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	resp := &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    statusCode,
			Message: msg,
		},
	}
	respByte, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error when marshal response to json, detail: ", err)
		return
	}
	_, err = w.Write(respByte)
	if err != nil {
		log.Printf("Error when write byte response, detail: ", err)
	}
}
