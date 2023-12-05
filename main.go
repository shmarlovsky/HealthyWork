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

// TODO: UI config

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
	app.Start(true)
	return app
}

var App = initApp()

func onReady() {
	systray.SetIcon(readIcon(hw.GetIcon(App.CurrentState, App.Running)))
	systray.SetTitle(hw.APPNAME)
	systray.SetTooltip(App.ShowStatus())

	continueMenuItem := systray.AddMenuItem("Continue", "Continue with current timer")
	continueMenuItem.Disable()
	pauseMenuItem := systray.AddMenuItem("Pause", "Pause timer")

	sitMenuItemText := fmt.Sprintf("Sit [%.0fmin]", hw.SecondsToMinutes(int(App.SitTime)))
	sitMenuItem := systray.AddMenuItem(sitMenuItemText, "Sit (timer reset)")

	standMenuItemText := fmt.Sprintf("Stand [%.0fmin]", hw.SecondsToMinutes(int(App.StandTime)))
	standMenuItem := systray.AddMenuItem(standMenuItemText, "Stand (timer reset)")

	systray.AddSeparator()
	quitMenuItem := systray.AddMenuItem("Quit", "Quit")

	for {
		select {
		case <-App.NotifyCh:
			systray.SetTooltip(App.ShowStatus())
		case <-continueMenuItem.ClickedCh:
			App.Continue()
			continueMenuItem.Disable()
			pauseMenuItem.Enable()
			systray.SetTooltip(App.ShowStatus())
			systray.SetIcon(readIcon(hw.GetIcon(App.CurrentState, App.Running)))
		case <-pauseMenuItem.ClickedCh:
			App.Pause()
			pauseMenuItem.Disable()
			continueMenuItem.Enable()
			systray.SetTooltip(App.ShowStatus())
			systray.SetIcon(readIcon(hw.GetIcon(App.CurrentState, App.Running)))
		case <-sitMenuItem.ClickedCh:
			App.DoSit(false)
			continueMenuItem.Disable()
			pauseMenuItem.Enable()
			systray.SetTooltip(App.ShowStatus())
			systray.SetIcon(readIcon(hw.GetIcon(App.CurrentState, App.Running)))
		case <-standMenuItem.ClickedCh:
			App.DoStand(false)
			continueMenuItem.Disable()
			pauseMenuItem.Enable()
			systray.SetTooltip(App.ShowStatus())
			systray.SetIcon(readIcon(hw.GetIcon(App.CurrentState, App.Running)))
		case <-quitMenuItem.ClickedCh:
			systray.Quit()
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
