package hw

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/toast.v1"
)

const (
	APPNAME                 = "HealthyWork"
	STAND_NOTIFY_TEXT       = "Time to stand!"
	SIT_NOTIFY_TEXT         = "Ok, sit down"
	UPDATE_TOOLTIP_INTERVAL = 5
)

type State string

const (
	SITTING  State = "sit"
	STANDING State = "stand"
)

type App struct {
	StandTime       int64
	SitTime         int64
	StartState      State
	CurrentState    State
	CurrentDuration int64
	// channel to send events (exchange info with systray e.g.)
	NotifyCh chan struct{}
}

func (a App) String() string {
	return fmt.Sprintf(
		"App(current: %v, time: %v; 	sit: %v, stand: %v, starting: %v)",
		a.CurrentState, a.CurrentDuration, a.SitTime, a.StandTime, a.StartState,
	)
}

func NewApp() *App {
	return &App{
		StandTime:       viper.GetInt64("standTime") * 60,
		SitTime:         viper.GetInt64("sitTime") * 60,
		StartState:      State(viper.GetString("startMode")),
		CurrentState:    State(viper.GetString("startMode")),
		CurrentDuration: 0,
		NotifyCh:        make(chan struct{}),
	}
}

func (a *App) DoSit(byTimer bool) {
	a.CurrentState = SITTING
	a.CurrentDuration = 0
	log.Printf("Go to %v", a.CurrentState)
	a.notify()
	if byTimer {
		showToast(SIT_NOTIFY_TEXT)
	}
}

func (a *App) DoStand(byTimer bool) {
	a.CurrentState = STANDING
	a.CurrentDuration = 0
	log.Printf("Go to %v", a.CurrentState)
	a.notify()
	if byTimer {
		showToast(STAND_NOTIFY_TEXT)
	}
}

func (a *App) SwitchState() {
	if a.CurrentState == SITTING {
		a.DoStand(true)
	} else if a.CurrentState == STANDING {
		a.DoSit(true)
	} else {
		panic("Unkown state. Should never got here")
	}
}

func (a *App) ShowStatus() string {
	minutes := SecondsToMinutes(int(a.CurrentDuration))
	s := string(a.CurrentState)
	duration := SecondsToMinutes(int(a.currentLimit()))
	percent := minutes / duration * 100
	return fmt.Sprintf("%ving: %.0f of %.0f minutes [%.0f%%]", strings.ToUpper(s[:1])+s[1:], minutes, duration, percent)
}

func (a *App) Reload() {
	new := NewApp()
	a.SitTime = new.SitTime
	a.StandTime = new.StandTime
	a.StartState = new.StartState
}

func (a *App) notify() {
	select {
	case a.NotifyCh <- struct{}{}:
	// in case no one waiting for the channel
	default:
	}
}

func (a *App) currentLimit() int64 {
	if a.CurrentState == SITTING {
		return a.SitTime
	} else if a.CurrentState == STANDING {
		return a.StandTime
	} else {
		panic("Unkown state. Should never got here")
	}
}

func (a *App) Go() {
	for {
		if a.CurrentDuration >= a.currentLimit() {
			a.SwitchState()
		}
		// send notify event every N of seconds
		if a.CurrentDuration%UPDATE_TOOLTIP_INTERVAL == 0 {
			a.notify()
		}
		// log.Printf("Working.. %v", a)
		time.Sleep(time.Second)
		a.CurrentDuration += 1
	}
}

func showToast(title string) {
	notification := toast.Notification{
		AppID: APPNAME,
		Title: title,
	}

	iconFullPath, err := filepath.Abs("assets/icon2.ico")
	if err == nil {
		notification.Icon = iconFullPath
	}

	err = notification.Push()
	if err != nil {
		log.Println(err)
	}
}

func SecondsToMinutes(seconds int) float64 {
	return float64(seconds) / 60
}
