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
	upgrader websocket.Upgrader
	count    atomic.Uint32
}

func (self *kmactor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	} else if !websocket.IsWebSocketUpgrade(r) {
		fmt.Fprintf(w, "kmactor %s @ %d", self.version, os.Getpid())
	} else if conn, err := self.upgrader.Upgrade(w, r, nil); err == nil {
		defer conn.Close()
		if self.count.CompareAndSwap(0, 1) {
			defer self.count.CompareAndSwap(1, 0)
			handled := 0
			count := 0
			cmd := Command{}
			log.Println("connected")
			defer func() { log.Printf("handled %d/%d\r\n", handled, count) }()
			for {
				cmd.Reset()
				if err = conn.ReadJSON(&cmd); err != nil {
					break
				} else if Play(&cmd) {
					handled += 1
				}
				count += 1
			}
		}
	}
}
