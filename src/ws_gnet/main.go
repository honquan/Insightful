package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/muster"
	"github.com/panjf2000/gnet/v2"
	repository "insightful/src/apis/repositories"
	"insightful/src/apis/services"
	services2 "insightful/src/ws_gnet/services"
	"log"
	_ "net/http/pprof"
	"time"
)

func init() {
	// Init services
	services.InitialServices()
}

func main() {
	var insightfullRepository repository.InsightfullRepository
	_ = services.GetServiceContainer().Invoke(func(s repository.InsightfullRepository) {
		insightfullRepository = s
	})

	var port int
	var multicore bool

	// Example command: go run main.go --port 8080 --multicore=true
	flag.IntVar(&port, "port", 8899, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()

	wss := &services2.WsServer{
		Addr:      fmt.Sprintf("tcp://localhost:%d", port),
		Multicore: multicore,
	}

	// init muster
	wss.Muster.MaxBatchSize = 1000
	wss.Muster.MaxConcurrentBatches = 5000
	wss.Muster.BatchTimeout = 5000 * time.Millisecond
	wss.Muster.PendingWorkCapacity = 3000
	wss.Muster.BatchMaker = func() muster.Batch {
		return &services2.WsServer{
			InsightfullRepo: insightfullRepository,
		}
	}
	err := wss.Muster.Start()
	if err != nil {
		panic(err)
	}

	// Start serving!
	log.Println("server exits:", gnet.Run(wss, wss.Addr, gnet.WithMulticore(multicore), gnet.WithReusePort(true), gnet.WithTicker(true)))
}
