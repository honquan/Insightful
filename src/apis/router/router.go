package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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
	a.Router.HandleFunc(fmt.Sprintf("%v/health-check", RouterWSPrefix), healthController.HealthCheck).Methods(http.MethodGet)

	// add router ws controller
	wsController := &controllers.WsController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/ws/worker", RouterWSPrefix), wsController.WsWorker).Methods(http.MethodPost)

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

	a.MysqlDB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}

	a.MysqlDB.DB().SetConnMaxLifetime(30 * time.Minute)

	// set max idle and max open conns
	a.MysqlDB.DB().SetMaxIdleConns(conf.EnvConfig.DBMysqlMaxIdleConns)
	a.MysqlDB.DB().SetMaxOpenConns(conf.EnvConfig.DBMysqlMaxOpenConns)
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
