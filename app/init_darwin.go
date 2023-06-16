//go:build darwin
// +build darwin

// macOS App 运行时的工作路径是 `/` 。
// 使用 `~/Library/Caches/kmactor` 作为工作路径，并将需要的资源文件复制到工作路径下。
//
// kmactor.app
// └── Contents
//        ├── Info.plist
//        ├── MacOS
//        │     └── kmactor
//        └── Resources
//               ├── config.json
//               └── icon.icns

package app

import (
	"io"
	"os"
	"path/filepath"
)

func Initialize() error {
	if path, err := os.Executable(); err != nil {
		return err
	} else if cacheDir, err := os.UserCacheDir(); err != nil {
		return err
	} else {
		cachePath := filepath.Join(cacheDir, filepath.Base(path))
		resPath := filepath.Join(path, "../..", "Resources")
		resOffset := len(resPath) + 1
		if err = os.MkdirAll(cachePath, 0750); err != nil {
			return err
		} else if err = os.Chdir(cachePath); err != nil {
			return err
		} else {
			return filepath.Walk(resPath, func(src string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				} else if info.IsDir() {
					return nil
				} else {
					res := src[resOffset:]
					if ignores[res] {
						return nil
					} else {
						dst := filepath.Join(cachePath, res)
						if _, err = os.Stat(dst); err == nil {
							return nil
						} else {
							return filecopy(dst, src)
						}
					}
				}
			})
		}
	}
}

var ignores = map[string]bool{
	"icon.icns": true,
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
