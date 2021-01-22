package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"insightful/controllers"
	"insightful/dtos"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	RouterWSPrefix = "/insightful"
)

type App struct {
	Router  *mux.Router
	MysqlDB *gorm.DB
}

func (a *App) InitRouter() {
	// creates a new instance of a mux router
	a.Router = mux.NewRouter().StrictSlash(true)

	// add router ws controller get
	wsController := &controllers.WsController{}
	a.Router.HandleFunc(fmt.Sprintf("%v/ws/get", RouterWSPrefix), wsController.WsGet).Methods(http.MethodGet)

	// register middleware
	a.Router.Use(a.recoverPanicMiddleware)
}

func (a *App) InitDatabase() {
	var err error
	mysqlUsername := os.Getenv("DB_MYSQL_USERNAME")
	mysqlPassword := os.Getenv("DB_MYSQL_PASSWORD")
	mysqlHost := os.Getenv("DB_MYSQL_HOST")
	mysqlPort := os.Getenv("DB_MYSQL_PORT")
	mysqlDBName := os.Getenv("DB_MYSQL_NAME")

	a.MysqlDB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}

	a.MysqlDB.DB().SetConnMaxLifetime(30 * time.Minute)

	// set max idle and max open conns
	//a.DB.DB().SetMaxIdleConns(os.Getenv("DB_MYSQL_MAXIDLECONNS"))
	//a.DB.DB().SetMaxOpenConns(os.Getenv("DB_MYSQL_MAXOPENCONNS"))
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
