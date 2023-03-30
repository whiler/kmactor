package main

import (
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/pkg/browser"

	"kmactor/app"
)

func tray(signals <-chan os.Signal, quit <-chan struct{}, status string, clean func()) {
	systray.Run(func() {
		systray.SetIcon(app.GetTrayIcon())
		systray.SetTooltip("kmactor")

		menuAbout := systray.AddMenuItem("关于", "About")
		menuStatus := systray.AddMenuItem("状态", "Status")
		menuQuit := systray.AddMenuItem("退出", "Quit")

		running := true
		for running {
			select {
			case <-menuAbout.ClickedCh:
				if err := browser.OpenURL("https://github.com/whiler/kmactor"); err != nil {
					log.Println(err)
				}
			case <-menuStatus.ClickedCh:
				if err := browser.OpenURL(status); err != nil {
					log.Println(err)
				}
			case <-menuQuit.ClickedCh:
				running = false
			case <-signals:
				running = false
			case <-quit:
				running = false
			}
		}

		close(menuQuit.ClickedCh)
		close(menuStatus.ClickedCh)
		close(menuAbout.ClickedCh)
		systray.Quit()
	}, clean)
}
