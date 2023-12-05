package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/shmarlovsky/HealthyWork/hw"
	"github.com/spf13/viper"
)

func readIcon(iconPath string) []byte {
	fullpath, err := filepath.Abs(iconPath)
	if err != nil {
		log.Fatalf("cannot find icon %vm %v", iconPath, err)
	}

	content, err := os.ReadFile(fullpath)
	if err != nil {
		log.Fatalf("error reading icon %vm %v", iconPath, err)
		return []byte{}
	}

	return content
}

func initApp() *hw.App {
	loadConfig()

	app := hw.NewApp()

	go func() {
		app.Go()
	}()

	return app
}

var App = initApp()

func onReady() {
	systray.SetIcon(readIcon("assets/icon2.ico"))
	systray.SetTitle(hw.APPNAME)
	systray.SetTooltip(App.ShowStatus())

	sitMenuItemText := fmt.Sprintf("Sit [%.0fmin]", hw.SecondsToMinutes(int(App.SitTime)))
	sitMenuItem := systray.AddMenuItem(sitMenuItemText, "Sit")

	standMenuItemText := fmt.Sprintf("Stand [%.0fmin]", hw.SecondsToMinutes(int(App.StandTime)))
	standMenuItem := systray.AddMenuItem(standMenuItemText, "Stand")

	systray.AddSeparator()
	quitMenuItem := systray.AddMenuItem("Quit", "Quit")

	for {
		select {
		case <-sitMenuItem.ClickedCh:
			App.DoSit(false)
		case <-standMenuItem.ClickedCh:
			App.DoStand(false)
		case <-quitMenuItem.ClickedCh:
			systray.Quit()
		case <-App.NotifyCh:
			systray.SetTooltip(App.ShowStatus())
		}
	}

}

func onExit() {
	// clean up here
}

func loadConfig() {
	viper.SetConfigName("conf")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}
}

func main() {
	systray.Run(onReady, onExit)
}
