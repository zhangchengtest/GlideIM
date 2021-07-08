package im

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WsServerOptions struct {
	Host          string
	Port          int
	ReadDeadLine  time.Duration
	WriteDeadLine time.Duration
}

type WsServer struct {
	options  *WsServerOptions
	upgrader websocket.Upgrader
}

// NewWsServer options can be nil, use default value when nil.
func NewWsServer(options *WsServerOptions) *WsServer {

	if options == nil {
		options = &WsServerOptions{
			Host:          "0.0.0.0",
			Port:          8080,
			ReadDeadLine:  12 * time.Minute,
			WriteDeadLine: 12 * time.Minute,
		}
	}
	ws := new(WsServer)
	ws.options = options
	ws.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 65536,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return ws
}

func (ws *WsServer) handleWebSocketRequest(writer http.ResponseWriter, request *http.Request) {

	conn, err := ws.upgrader.Upgrade(writer, request, nil)
	if err != nil {

		return
	}

	con := NewWsConnection(conn, ws.options)
	NewClient(con).Run()
}

func (ws *WsServer) Start() {

	http.HandleFunc("/ws", ws.handleWebSocketRequest)

	addr := fmt.Sprintf("%s:%d", ws.options.Host, ws.options.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}

}
