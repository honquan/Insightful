package job_worker

import (
	"log"
)

// Job represents the job to be run
type Job struct {
	Payload interface{}
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job

// Callback function fire after recieve Job
type WorkerInstanceFunc func(job Job) error

// Worker represents the worker that executes the job
type Worker struct {
	WorkerFunc WorkerInstanceFunc
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job, WorkerFunc WorkerInstanceFunc) Worker {
	return Worker{
		WorkerFunc: WorkerFunc,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				if w.WorkerFunc == nil {
					log.Printf("Callback worker func can not be nil")
					continue
				}

				// we have received a work request.
				if err := w.WorkerFunc(job); err != nil {
					log.Printf("Error when fire worker: %s", err.Error())
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
