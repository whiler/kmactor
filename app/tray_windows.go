//go:build windows
// +build windows

package app

import _ "embed"

//go:embed black.ico
var trayIcon []byte

func GetTrayIcon() []byte { return trayIcon }
