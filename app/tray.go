package app

import _ "embed"

//go:embed tray.ico
var trayIcon []byte

func GetTrayIcon() []byte { return trayIcon }
