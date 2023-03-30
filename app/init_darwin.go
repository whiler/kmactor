//go:build darwin
// +build darwin

// macOS App 运行时的工作路径是 `/` 。
// 使用 `~/Library/Caches/kmactor` 作为工作路径，并将需要的资源文件复制到工作路径下。

package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

var files = []string{
	"cert.pem",
	"key.pem",
	"repo.txt",
}

func Initialize(cli bool) error {
	if cli {
		return nil
	} else if path, err := os.Executable(); err != nil {
		return err
	} else if cacheDir, err := os.UserCacheDir(); err != nil {
		return err
	} else {
		contentsPath := filepath.Join(path, "../..")
		base := filepath.Base(filepath.Join(contentsPath, ".."))
		name := base[:strings.LastIndex(base, ".app")]
		cachePath := filepath.Join(cacheDir, name)
		if err = os.MkdirAll(cachePath, 0750); err != nil {
			return err
		} else if err = os.Chdir(cachePath); err != nil {
			return err
		} else {
			for _, cur := range files {
				src := filepath.Join(contentsPath, "Resources", cur)
				dst := filepath.Join(cachePath, cur)
				if !exists(src) || exists(dst) {
					continue
				} else if err = filecopy(dst, src); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func filecopy(dst, src string) error {
	if r, err := os.Open(src); err != nil {
		return err
	} else {
		defer r.Close()
		if w, err := os.Create(dst); err != nil {
			return err
		} else {
			defer w.Close()
			_, err = io.Copy(w, r)
			return err
		}
	}
}
