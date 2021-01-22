package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

type BaseController struct {}

func (b *BaseController) ServeJSONWithCode(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
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

func (b *BaseController) ServeBytesWithCode(w http.ResponseWriter, statusCode int, respByte []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(respByte)
	if err != nil {
		log.Printf("Error when write byte response, detail: ", err)
	}
}
