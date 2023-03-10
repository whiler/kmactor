package main

import (
	"context"
	"errors"
	"flag"
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
)

var (
	normal        = "0.1.2"
	preRelease    = "dev"
	buildRevision string

	ver string
)

func use(path string) string {
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		return ""
	} else {
		return path
	}
}

func main() {
	var (
		flagSet   = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		port      int
		token     string
		cert, key string
		logto     string
		repo      string
		version   bool
	)

	log.SetFlags(log.Ltime)

	flagSet.IntVar(&port, "port", 9242, "local port")
	flagSet.StringVar(&token, "token", "", "token")
	flagSet.StringVar(&cert, "cert", use("cert.pem"), "cert file path")
	flagSet.StringVar(&key, "key", use("key.pem"), "key file path")
	flagSet.StringVar(&logto, "log", "kmactor.log", "log file path")
	flagSet.StringVar(&repo, "repo", use("repo.txt"), "auto update cert from repo")
	flagSet.BoolVar(&version, "version", false, "version")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Println(err)
	} else if version {
		fmt.Println(ver)
	} else if closer, err := logging(logto); err != nil {
		log.Println(err)
	} else {
		defer closer.Close()
		if 0 >= port || port >= 65536 {
			log.Printf("invalid port: %d", port)
		} else if (cert != "" && key == "") || (cert == "" && key != "") {
			log.Println("cert and key are required at the same time")
		} else if err = updateCert(cert, key, repo); err != nil {
			log.Println(err)
		} else if names, err := getCertNames(cert, key); err != nil {
			log.Println(err)
		} else if handler, err := kmactor.Build(ver, token); err != nil {
			log.Println(err)
		} else {
			tls := cert != "" && key != ""
			srv := &http.Server{
				Addr:    fmt.Sprintf("localhost:%d", port),
				Handler: h2c.NewHandler(handler, &http2.Server{}),
			}

			quit := make(chan struct{})
			go func() {
				defer close(quit)
				var err error
				if tls {
					err = srv.ListenAndServeTLS(cert, key)
				} else {
					err = srv.ListenAndServe()
				}
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Printf("serve error: %v", err)
				}
			}()

			select {
			case <-quit:
			case <-time.After(time.Second):
				proto := "ws"
				if tls {
					proto = "wss"
				}
				for _, name := range names {
					log.Printf("serving at %s://%s:%d", proto, name, port)
				}

				sig := make(chan os.Signal, 1)
				signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
				select {
				case <-sig:
				case <-quit:
				}
				signal.Stop(sig)
				close(sig)

				ctx, cancel := context.WithTimeout(context.Background(), 13*time.Second)
				if err := srv.Shutdown(ctx); err != nil {
					log.Printf("shutdown error: %v", err)
				}
				cancel()
			}
		}
	}
}

type dummp struct{}

func (*dummp) Close() error { return nil }

func logging(path string) (io.Closer, error) {
	if path == "-" {
		return &dummp{}, nil
	} else if file, err := os.Create(path); err != nil {
		return nil, err
	} else {
		log.SetOutput(io.MultiWriter(log.Writer(), file))
		return file, nil
	}
}

func init() {
	ver = fmt.Sprintf("%s-%s+%s", normal, preRelease, buildRevision)
}
