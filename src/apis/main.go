package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/evanphx/wildcat"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/fasthttp"
	"insightful/src/apis/dtos"
	"insightful/src/apis/router"
	"insightful/src/apis/services"
	"log"
	"net/http"
	//_ "net/http/pprof"
	"time"
)

func init() {
	// Init logging
	//config := logger.Configuration{
	//	EnableConsole:     true,
	//	ConsoleLevel:      strings.ToLower(conf.EnvConfig.LogLevel),
	//	ConsoleJSONFormat: true,
	//	EnableFile:        false,
	//}
	//logger := logger.NewLogger(config, logger.InstanceZapLogger)
	//if logger == nil {
	//	log.Printf("Could not instantiate log")
	//}

	// Init services
	services.InitialServices()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// run custom job worker without redis
	//custom_worker.JobQueue = make(chan custom_worker.Job, conf.EnvConfig.MaxWorker)
	//dispatcher := custom_worker.NewDispatcher(conf.EnvConfig.MaxWorker)
	//dispatcher.Run()

	// run muster
	//muster.Run()

	// init router
	a := router.App{}
	a.InitRouter()

	// run worker go worker
	//worker.RunGoWorker()

	// run go craft
	//go worker.RunGoCraft()

	//wsController := controllers.NewWebsocketController()
	//http.HandleFunc("/test/ws/worker-ants", wsController.WebsocketAntsWorker)

	// run
	a.Run(":8899")
}

func main_fast_http() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	myHandler := &MyHandler{}
	fasthttp.ListenAndServe(":8899", myHandler.HandleFastHTTP)
}

// ////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////   FASTHTTP   ////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////
type MyHandler struct {
	foobar string
}

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

// request handler in net/http style, i.e. method bound to MyHandler struct.
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	start := time.Now()

	obj := &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "Ok",
		},
	}
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		elapsed := time.Since(start)
		fmt.Print("", elapsed, err.Error())
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////// GNET //////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////
var (
	errMsg      = "Internal Server Error"
	errMsgBytes = []byte(errMsg)
)

type httpServer struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
}

type httpCodec struct {
	parser *wildcat.HTTPParser
	buf    []byte
}

func (hc *httpCodec) appendResponse() {
	resp := &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code: http.StatusOK,
		},
	}

	respByte, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error when marshal response to json, detail: ", err)
		return
	}
	fmt.Print(respByte)

	hc.buf = append(hc.buf, "HTTP/1.1 200 OK\r\nContent-Type: application/json"...)
	hc.buf = append(hc.buf, "Hello World!"...)
}

func (hs *httpServer) OnBoot(eng gnet.Engine) gnet.Action {
	hs.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n", hs.multicore, hs.addr)
	return gnet.None
}

func (hs *httpServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(&httpCodec{parser: wildcat.NewHTTPParser()})
	return nil, gnet.None
}

func (hs *httpServer) OnTraffic(c gnet.Conn) gnet.Action {
	hc := c.Context().(*httpCodec)
	buf, _ := c.Next(-1)

pipeline:
	headerOffset, err := hc.parser.Parse(buf)
	if err != nil {
		c.Write(errMsgBytes)
		return gnet.Close
	}
	hc.appendResponse()
	bodyLen := int(hc.parser.ContentLength())
	if bodyLen == -1 {
		bodyLen = 0
	}
	buf = buf[headerOffset+bodyLen:]
	if len(buf) > 0 {
		goto pipeline
	}

	c.Write(hc.buf)
	hc.buf = hc.buf[:0]
	return gnet.None
}

func main_gnet() {
	var port int
	var multicore bool

	// Example command: go run main.go --port 8080 --multicore=true
	flag.IntVar(&port, "port", 9080, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()

	hs := &httpServer{addr: fmt.Sprintf("tcp://127.0.0.1:%d", port), multicore: multicore}

	// Start serving!
	log.Println("server exits:", gnet.Run(hs, hs.addr, gnet.WithMulticore(multicore)))
}
