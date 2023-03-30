//go:build !windows
// +build !windows

package app

import _ "embed"

//go:embed 96.png
var trayIcon []byte

func GetTrayIcon() []byte { return trayIcon }
