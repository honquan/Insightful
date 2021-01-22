package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"insightful/controllers"
	"insightful/dtos"
	"log"
	"net/http"
)

const (
	RouterWSPrefix = "/insightful"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
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
