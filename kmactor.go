package kmactor

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type kmactor struct {
	version   string
	token     string
	upgrader  websocket.Upgrader
	count     atomic.Uint32
	birthtime time.Time
}

func (self *kmactor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	} else if !websocket.IsWebSocketUpgrade(r) {
		fmt.Fprintf(w, "app: kmactor\r\n")
		fmt.Fprintf(w, "version: %s\r\n", self.version)
		fmt.Fprintf(w, "age: %s\r\n", time.Since(self.birthtime))
		fmt.Fprintf(w, "process id: %d\r\n", os.Getpid())
		fmt.Fprintf(w, "session count: %d\r\n", self.count.Load())
		scheme := "ws"
		if r.TLS != nil {
			scheme = "wss"
		}
		fmt.Fprintf(w, "address: %s\r\n", (&url.URL{Scheme: scheme, Host: r.Host}).String())
		if len(self.token) > 0 {
			fmt.Fprintf(w, "token: %s\r\n", self.token)
		}
	} else if len(self.token) > 0 && self.token != r.URL.Query().Get("token") {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	} else if conn, err := self.upgrader.Upgrade(w, r, nil); err == nil {
		defer conn.Close()
		if self.count.CompareAndSwap(0, 1) {
			defer self.count.CompareAndSwap(1, 0)
			handled := 0
			count := 0
			width, height := GetScreenSize()
			cmd := Command{}
			log.Println("connected")
			defer func() { log.Printf("handled %d/%d", handled, count) }()
			for {
				cmd.Reset()
				if err = conn.ReadJSON(&cmd); err != nil {
					break
				} else if Play(&cmd, width, height) {
					handled += 1
				}
				count += 1
			}
		} else {
			log.Println("refused")
		}
	}
}

func Build(ver, token string) (http.Handler, error) {
	return &kmactor{
		version: ver,
		token:   token,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(*http.Request) bool { return true },
		},
		birthtime: time.Now(),
	}, nil
}
