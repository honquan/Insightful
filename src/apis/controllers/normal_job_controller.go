package controllers

import (
	"encoding/json"
	"github.com/gocraft/work"
	"insightful/src/apis/dtos"
	go_worker "insightful/src/apis/pkg/worker"
	"io"
	"net/http"
	"time"
)

type NormalJobController struct {
	BaseController
}

func (s *NormalJobController) NormalJobWorkerGoWorker(w http.ResponseWriter, r *http.Request) {
	//parse
	//var req dtos.WsPayload
	//if err := schema.NewDecoder().Decode(&req, r.URL.Query()); err != nil {
	//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	_ = go_worker.AddJob("Sample", time.Now().UTC(), "arg1")

	s.ServeJSONWithCode(w, http.StatusOK, &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "Ok",
		},
	})
}

func (s *NormalJobController) NormalJobWorkerGoCraft(w http.ResponseWriter, r *http.Request) {
	//parse
	var content = &dtos.WsPayload{}
	err := json.NewDecoder(io.LimitReader(r.Body, 2048)).Decode(&content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//
	enqueueJobCraft(
		"send_email",
		work.Q{"data": content},
	)

	s.ServeJSONWithCode(w, http.StatusOK, &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "Ok",
		},
	})
}
