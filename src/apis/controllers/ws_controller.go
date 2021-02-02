package controllers

import (
	"encoding/json"
	"fmt"
	"insightful/src/apis/conf"
	"insightful/src/apis/dtos"
	worker "insightful/src/apis/kit/job_worker"
	"io"
	"log"
	"net/http"
)

type WsController struct {
	BaseController
}

func (s *WsController) WsWorker(w http.ResponseWriter, r *http.Request) {
	log.Println("Start ws worker")
	d := worker.NewDispatcher(conf.EnvConfig.MaxWorker).AppendCallbackWorker(s.FireWorker).Run()

	// Read the body into a string for json decoding
	var content = &dtos.WsPayload{}
	err := json.NewDecoder(io.LimitReader(r.Body, 2048)).Decode(&content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Push the work onto the queue.
	d.Submit(*content)

	s.ServeJSONWithCode(w, http.StatusOK, &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "Ok",
		},
	})
}

func (s *WsController) FireWorker(job worker.Job) error {
	fmt.Printf("%+v\n", job)
	return nil
}
