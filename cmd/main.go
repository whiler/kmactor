package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"kmactor"
)

var (
	normal        = "0.1.1"
	preRelease    = "dev"
	buildRevision string

	ErrNoCertificate   = errors.New("no cert")
	ErrTimeCertificate = errors.New("out of cert time")
)

type dummp struct{}

func (*dummp) Close() error { return nil }

func main() {
	var (
		flagSet   = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		port      int
		token     string
		cert, key string
		logto     string
		version   bool
	)

	log.SetFlags(log.Ltime)

	flagSet.IntVar(&port, "port", 9242, "local port")
	flagSet.StringVar(&token, "token", "", "token")
	flagSet.StringVar(&cert, "cert", ensure("cert.pem"), "cert file path")
	flagSet.StringVar(&key, "key", ensure("key.pem"), "key file path")
	flagSet.StringVar(&logto, "log", "kmactor.log", "log file path")
	flagSet.BoolVar(&version, "version", false, "version")

	ver := fmt.Sprintf("%s-%s+%s", normal, preRelease, buildRevision)

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
		} else if names, err := getCertName(cert, key); err != nil {
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

func ensure(path string) string {
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		return ""
	} else {
		return path
	}
}

func getCertName(certpath, keypath string) ([]string, error) {
	if certpath == "" || keypath == "" {
		return []string{"localhost"}, nil
	} else {
		cur := time.Now()
		if cert, err := tls.LoadX509KeyPair(certpath, keypath); err != nil {
			return nil, err
		} else if len(cert.Certificate) == 0 {
			return nil, ErrNoCertificate
		} else if x509cert, err := x509.ParseCertificate(cert.Certificate[0]); err != nil {
			return nil, err
		} else if cur.Before(x509cert.NotBefore) || cur.After(x509cert.NotAfter) {
			return nil, ErrTimeCertificate
		} else {
			set := map[string]bool{}
			set[x509cert.Subject.CommonName] = true
			for _, name := range x509cert.DNSNames {
				set[name] = true
			}
			for _, ip := range x509cert.IPAddresses {
				set[ip.String()] = true
			}
			names := make([]string, 0, len(set))
			for name := range set {
				names = append(names, strings.ReplaceAll(name, "*", "local"))
			}
			return names, nil
		}
	}
}
