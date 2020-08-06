package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kordar/ws/iface"
	"github.com/kordar/ws/net"
	"github.com/kordar/ws/utils"
	"net/http"
)

const (
	LinkMaxUint32        = ^uint32(0)
	DefaultUuidLinkCache = 1000
	RoomMaxUint32        = 1000
	DefaultUuidRoomCache = 50
)

var (
	connUUID = utils.NewUUIDGenerator(LinkMaxUint32, DefaultUuidLinkCache)
	roomUUID = utils.NewUUIDGenerator(RoomMaxUint32, DefaultUuidRoomCache)
	roomMgr = net.NewRoomManager()
	room = net.NewRoom(roomUUID.GetUint32())
)

func init()  {
	roomMgr.GetMsgHandler().StartWorkerPool()
}

type PingRouter struct {
	net.BaseRouter
}

func (r *PingRouter) Handle(request iface.IRequest) {
	conn := request.GetConnection()
	conn.GetRoom().Broadcast(request.GetMsgID(), request.GetMessage(), request.GetData(),conn.GetConnID())
	// _ = conn.SendMsg(request.GetMsgID(), request.GetMessage(), request.GetData())
}

// WsPage is a websocket handler
func WsPage(c *gin.Context) {

	// change the request to websocket model
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	if err := roomMgr.JoinRoom(room, 2); err != nil {
		fmt.Println("join, ", err)
		_ = conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		return
	}

	connID := connUUID.GetUint32()
	dealConn := net.NewConnection(room, conn, connID, roomMgr.GetMsgHandler())
	go dealConn.Start()
}

func main() {

	roomMgr.AddRouter(1, &PingRouter{})

	router := gin.Default()

	router.StaticFS("/html", http.Dir("./html"))

	router.GET("/ws", WsPage)
	_ = router.Run(":8000")
}
