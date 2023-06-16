package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var errTimeCertificate = errors.New("out of cert time")

func validate(cert tls.Certificate) error {
	cur := time.Now()
	if x509cert, err := x509.ParseCertificate(cert.Certificate[0]); err != nil {
		return err
	} else if cur.Before(x509cert.NotBefore) || cur.After(x509cert.NotAfter) {
		return errTimeCertificate
	} else {
		return nil
	}
}

func isValid(certpath, keypath string) bool {
	if cert, err := tls.LoadX509KeyPair(certpath, keypath); err == nil {
		return validate(cert) == nil
	} else {
		return false
	}
}

func getCertNames(certpath, keypath string) ([]string, error) {
	if certpath == "" || keypath == "" {
		return []string{"localhost"}, nil
	} else if cert, err := tls.LoadX509KeyPair(certpath, keypath); err != nil {
		return nil, err
	} else if x509cert, err := x509.ParseCertificate(cert.Certificate[0]); err != nil {
		return nil, err
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
		sort.Slice(names, func(i, j int) bool { return len(names[i]) < len(names[j]) })
		return names, nil
	}
}

func wget(target, version string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil); err != nil {
		return nil, err
	} else {
		req.Header.Add("User-Agent", "kmactor/"+version)
		if resp, err := http.DefaultClient.Do(req); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()
			return io.ReadAll(resp.Body)
		}
	}
}

func dumpto(path string, content []byte) error {
	if file, err := os.CreateTemp("", filepath.Base(path)); err != nil {
		return err
	} else {
		temp := file.Name()
		if wrote, err := file.Write(content); err != nil {
			file.Close()
			os.Remove(temp)
			return err
		} else if wrote != len(content) {
			file.Close()
			os.Remove(temp)
			return io.ErrShortWrite
		} else {
			file.Close()
			return os.Rename(temp, path)
		}
	}
}

func fetchCert(repo, version string) ([]byte, []byte, error) {
	log.Printf("fetching cert from %s", repo)
	if u, err := url.Parse(repo); err != nil {
		return nil, nil, err
	} else if certContent, err := wget(u.JoinPath("cert.pem").String(), version); err != nil {
		return nil, nil, err
	} else if keypath, err := wget(u.JoinPath("key.pem").String(), version); err != nil {
		return nil, nil, err
	} else {
		return certContent, keypath, nil
	}
}

func updateCert(certpath, keypath, repo, version string) error {
	if repo == "" || certpath == "" || keypath == "" {
		return nil
	} else if isValid(certpath, keypath) {
		return nil
	} else if certContent, keyContent, err := fetchCert(repo, version); err != nil {
		return err
	} else if cert, err := tls.X509KeyPair(certContent, keyContent); err != nil {
		return err
	} else if err = validate(cert); err != nil {
		return err
	} else if err = dumpto(certpath, certContent); err != nil {
		return err
	} else if err = dumpto(keypath, keyContent); err != nil {
		return err
	} else {
		return nil
	}
}
