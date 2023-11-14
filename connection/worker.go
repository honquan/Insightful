package connection

import (
	"insightful/common/config"
	"insightful/src/apis/kit/custom_worker"
)

func InitWorker(conf *config.Config) *custom_worker.Dispatcher {

	custom_worker.JobQueue = make(chan custom_worker.Job, conf.Worker.MaxWorker)
	return custom_worker.NewDispatcher(conf.Worker.MaxWorker)
}
