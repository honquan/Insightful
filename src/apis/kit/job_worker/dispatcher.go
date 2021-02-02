package job_worker

import (
	"fmt"
	"insightful/src/apis/conf"
	"log"
)

type Dispatcher struct {
	WorkerInc WorkerInstanceFunc
	// A pool of workers channels that are registered with the dispatcher
	maxWorkers int
	WorkerPool chan chan Job
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) AppendCallbackWorker(workerInc WorkerInstanceFunc) *Dispatcher {
	d.WorkerInc = workerInc
	return d
}

func (d *Dispatcher) Run() *Dispatcher {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.WorkerInc)
		worker.Start()
	}
	go d.dispatch()

	return d
}

func (d *Dispatcher) dispatch() {
	fmt.Println("Worker que dispatcher started...")
	for {

		select {
		case job := <-JobQueue:
			log.Printf("a dispatcher request received")
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

func (d *Dispatcher) Submit(data interface{}) {
	JobQueue = make(chan Job, conf.EnvConfig.MaxQueue)
	JobQueue <- Job{Payload: data}
}
