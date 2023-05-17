package controllers

import (
	"fmt"
	"github.com/gocraft/work"
	"insightful/src/apis/dtos"
	go_worker "insightful/src/apis/pkg/worker"
	"log"
	"net/http"
	"time"
)

type JobController struct {
	BaseController
}

func (s *JobController) WorkerGoWorker(w http.ResponseWriter, r *http.Request) {
	workerTime := time.Now().UTC()
	jId := go_worker.AddJob("Sample", workerTime, "arg1", "arg2")
	jId2 := go_worker.AddJob("Sample2", workerTime, "arg3", "arg4")
	fmt.Println("jid: %s", jId)
	fmt.Printf("Working on job, arg: %s", jId2)

	s.ServeJSONWithCode(w, http.StatusOK, &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "Ok",
		},
	})
}

func (s *JobController) WorkerGoCraft(w http.ResponseWriter, r *http.Request) {
	enqueueEmail()

	s.ServeJSONWithCode(w, http.StatusOK, &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "Ok",
		},
	})
}

var enqueuer = work.NewEnqueuer("my_app_namespace", go_worker.RedisPool)

func enqueueJobCraft(job string, payload work.Q) {
	_, err := enqueuer.Enqueue(job, payload)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Enqueued:", job, ", with Payload:", payload)
}

func enqueueEmail() {
	enqueueJobCraft(
		"send_email",
		work.Q{"address": "test@example.com", "subject": "hello world", "customer_id": 4},
	)
}
