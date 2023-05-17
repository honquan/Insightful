package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"insightful/src/apis/conf"
	"insightful/src/apis/controllers"
	"insightful/src/apis/dtos"
	"log"
	"net/http"
	"time"
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

	// custom job worker
	customJobController := &controllers.CustomJobController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/custom/worker", RouterWSPrefix), customJobController.WorkerJobCustom).Methods(http.MethodPost)

	// add router websocket controller
	wsController := &controllers.WsController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/ws/worker", RouterWSPrefix), wsController.WebsocketWorker).Methods(http.MethodGet)

	// normal job
	normalController := &controllers.NormalJobController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/job-normal/go-worker", RouterWSPrefix), normalController.NormalJobWorkerGoWorker).Methods(http.MethodPost)
	a.Router.HandleFunc(fmt.Sprintf("%v/job-normal/craft-worker", RouterWSPrefix), normalController.NormalJobWorkerGoCraft).Methods(http.MethodPost)

	// job
	jobController := &controllers.JobController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/job/worker-goworker", RouterWSPrefix), jobController.WorkerGoWorker).Methods(http.MethodGet)
	a.Router.HandleFunc(fmt.Sprintf("%v/job/worker-gocraft", RouterWSPrefix), jobController.WorkerGoCraft).Methods(http.MethodGet)

	// register middleware
	a.Router.Use(a.recoverPanicMiddleware)
}

func (a *App) InitDatabase() {
	var err error
	mysqlUsername := conf.EnvConfig.DBMysqlUsername
	mysqlPassword := conf.EnvConfig.DBMysqlPassword
	mysqlHost := conf.EnvConfig.DBMysqlHost
	mysqlPort := conf.EnvConfig.DBMysqlPort
	mysqlDBName := conf.EnvConfig.DBMysqlName

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)
	a.MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}

	sqlDB, err := a.MysqlDB.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// set max idle and max open conns
	sqlDB.SetMaxIdleConns(conf.EnvConfig.DBMysqlMaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.EnvConfig.DBMysqlMaxOpenConns)
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
