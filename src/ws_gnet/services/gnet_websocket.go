package services

import (
	"context"
	"fmt"
	"github.com/facebookgo/muster"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"insightful/model"
	repository "insightful/src/apis/repositories"
	"insightful/src/ws_gnet/kit/gnet_worker"
	"os"
	"sync/atomic"
	"time"
	"unsafe"
)

type WsServer struct {
	gnet.BuiltinEventEngine

	Addr      string
	Multicore bool
	eng       gnet.Engine
	connected int64

	InsightfullRepo repository.InsightfullRepository
	Muster          muster.Client
	Items           []mongo.WriteModel
}

func (wss *WsServer) OnBoot(eng gnet.Engine) gnet.Action {
	wss.eng = eng
	logging.Infof("echo server with multi-core=%t is listening on %s", wss.Multicore, wss.Addr)
	return gnet.None
}

func (wss *WsServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(new(gnet_worker.WsCodec))
	atomic.AddInt64(&wss.connected, 1)
	return nil, gnet.None
}

func (wss *WsServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		//logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt64(&wss.connected, -1)
	//logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}

func (wss *WsServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ws := c.Context().(*gnet_worker.WsCodec)
	if ws.ReadBufferBytes(c) == gnet.Close {
		return gnet.Close
	}
	ok, action := ws.Upgrade(c)
	if !ok {
		return
	}

	if ws.Buf.Len() <= 0 {
		return gnet.None
	}
	messages, err := ws.Decode(c)
	if err != nil {
		return gnet.Close
	}
	if messages == nil {
		return
	}
	for _, message := range messages {
		//logging.Infof("conn[%v] receive [op=%v] [msg=%v..., len=%d]", c.RemoteAddr().String(), message.OpCode, string(message.Payload[:128]), len(message.Payload))

		//var data interface{}
		//_ = json.Unmarshal(message.Payload, &data)
		wss.Push(model.Insightful{
			Mongo: model.Mongo{
				CreatedAt: time.Now().Unix(),
				UpdatedAt: 0,
			},
			Coordinates: ByteSlice2String(message.Payload),
		})

		// This is the echo server
		/*err = wsutil.WriteServerMessage(c, message.OpCode, message.Payload)
		if err != nil {
			//logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
			return gnet.Close
		}*/
	}
	return gnet.None
}

func (wss *WsServer) OnTick() (delay time.Duration, action gnet.Action) {
	//logging.Infof("[connected-count=%v]", atomic.LoadInt64(&wss.connected))
	return 3 * time.Second, gnet.None
}

func (s *WsServer) Stop() error {
	return s.Muster.Stop()
}

// The CoordinateClient provides a typed Add method which enqueues the work.
func (s *WsServer) Push(item interface{}) {
	s.Muster.Work <- item
}

// The batch provides an untyped Add to satisfy the muster.Batch interface. As
// is the case here, the Batch implementation is internal to the user of muster
// and not exposed to the users of ShoppingClient.
func (s *WsServer) Add(item interface{}) {
	s.Items = append(s.Items, mongo.NewInsertOneModel().SetDocument(item))
}

// Once a Batch is ready, it will be Fired. It must call notifier.Done once the
// batch has been processed.
func (s *WsServer) Fire(notifier muster.Notifier) {
	defer notifier.Done()

	//err := s.InsightfullRepo.CreateMany(context.Background(), s.Items)
	err := s.InsightfullRepo.BulkWrite(context.Background(), s.Items)
	if err != nil {
		fmt.Println("error when create many mongo:", err)
	}
	os.Stdout.Sync()
}

func ByteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func convert(myBytes []byte) string {
	//return unsafe.String((unsafe.SliceData(myBytes)), len(myBytes))
	return string(myBytes[:])
}
