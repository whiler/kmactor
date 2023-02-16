package kmactor

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type kmactor struct {
	version  string
	token    string
	upgrader websocket.Upgrader
	count    atomic.Uint32
}

func (self *kmactor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	} else if !websocket.IsWebSocketUpgrade(r) {
		fmt.Fprintf(w, "kmactor %s @ %d", self.version, os.Getpid())
	} else if self.token != r.URL.Query().Get("token") {
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
