package controllers

import (
	"encoding/json"
	"fmt"
	"insightful/src/apis/dtos"
	jobWorker "insightful/src/apis/kit/custom_worker"
	"io"
	"log"
	"net/http"
)

type CustomJobController struct {
	BaseController
}

func (s *CustomJobController) WorkerJobCustom(w http.ResponseWriter, r *http.Request) {
	log.Println("Receive request")
	// Read the body into a string for json decoding
	var content = &dtos.WsPayload{}
	err := json.NewDecoder(io.LimitReader(r.Body, 2048)).Decode(&content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Push the work onto the queue.
	//work := jobWorker.Job{Payload: content}
	//jobWorker.JobQueue <- work

	// Push the work onto the queue.
	jobWorker.Submit(*content)

	s.ServeJSONWithCode(w, http.StatusOK, &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "Ok",
		},
	})
}

func (s *CustomJobController) FireWorker(job jobWorker.Job) error {
	fmt.Printf("%+v\n", job)
	return nil
}
