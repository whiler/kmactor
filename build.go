package kmactor

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func Build(ver, token string) (http.Handler, error) {
	return &kmactor{
		version: ver,
		token:   token,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(*http.Request) bool { return true },
		},
	}, nil
}
