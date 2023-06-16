package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"kmactor"
	"kmactor/app"
)

var (
	normal        = "0.2.0"
	preRelease    = "dev"
	buildRevision = "197001.0000000"

	ver = fmt.Sprintf("%s-%s+%s", normal, preRelease, buildRevision)
)

func main() {
	log.SetFlags(log.Ltime)
	if err := app.Initialize(); err != nil {
		log.Println(err.Error())
	} else if cfg, err := loadConfig("config.json", 9242, "kmactor.log"); err != nil {
		log.Println(err.Error())
	} else if closer, err := logto(cfg.Log); err != nil {
		log.Println(err.Error())
	} else if 0 >= cfg.Port || cfg.Port >= 65536 {
		log.Printf("invalid port: %d", cfg.Port)
		closer.Close()
	} else if (cfg.Cert != "" && cfg.Key == "") || (cfg.Cert == "" && cfg.Key != "") {
		log.Println("cert and key are required at the same time")
		closer.Close()
	} else if err = updateCert(cfg.Cert, cfg.Key, cfg.Repo, ver); err != nil {
		log.Println(err.Error())
		closer.Close()
	} else if names, err := getCertNames(cfg.Cert, cfg.Key); err != nil {
		log.Println(err.Error())
		closer.Close()
	} else if handler, err := kmactor.Build(ver, cfg.Token); err != nil {
		log.Println(err.Error())
		closer.Close()
	} else {
		tls := cfg.Cert != "" && cfg.Key != ""
		srv := &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", cfg.Port),
			Handler: h2c.NewHandler(handler, &http2.Server{}),
		}

		quit := make(chan struct{})
		go func() {
			defer close(quit)
			var e error
			if tls {
				e = srv.ListenAndServeTLS(cfg.Cert, cfg.Key)
			} else {
				e = srv.ListenAndServe()
			}
			if e != nil && !errors.Is(e, http.ErrServerClosed) {
				log.Printf("serve error: %s", e.Error())
			}
		}()

		select {
		case <-quit:
			closer.Close()
		case <-time.After(time.Second):
			ws := "ws"
			proto := "http"
			if tls {
				ws = "wss"
				proto = "https"
			}
			for _, name := range names {
				log.Printf("serving at %s://%s:%d", ws, name, cfg.Port)
			}

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			clean := func() {
				signal.Stop(sig)
				close(sig)
				ctx, cancel := context.WithTimeout(context.Background(), 13*time.Second)
				if e := srv.Shutdown(ctx); e != nil {
					log.Printf("shutdown error: %s", e.Error())
				}
				cancel()
				log.Println("quit")
				closer.Close()
			}
			tray(sig, quit, fmt.Sprintf("%s://%s:%d", proto, names[0], cfg.Port), clean)
		}
	}
}

func logto(path string) (io.Closer, error) {
	if path == "-" {
		return &dummp{}, nil
	} else if file, err := os.Create(path); err != nil {
		return nil, err
	} else {
		log.SetOutput(file)
		return file, nil
	}
}

type dummp struct{}

func (*dummp) Close() error {
	return nil
}
