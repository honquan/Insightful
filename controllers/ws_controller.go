package controllers

import "net/http"

type WsController struct {
}

func (s *WsController) WsGet(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write([]byte("Hello world"))
}
