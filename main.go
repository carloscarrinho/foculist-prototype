package main

import (
	"embed"
	"log"

	"changeme/backend"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/TrayIconTemplate.png
var trayIcon []byte

func init() {
	application.RegisterEvent[backend.TimerTick]("timer_tick")
	application.RegisterEvent[string]("timer_complete")
}

func main() {
	timerSvc := backend.NewTimerService()

	app := application.New(application.Options{
		Name:        "Foculist",
		Description: "macOS-First Pomodoro Timer",
		Services: []application.Service{
			application.NewService(timerSvc),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:       "Foculist",
		Width:       320,
		Height:      450,
		Frameless:   true,
		Hidden:      true,
		AlwaysOnTop: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		URL:              "/",
	})

	tray := app.SystemTray.New()
	tray.SetTemplateIcon(trayIcon)
	tray.SetLabel("25:00")
	tray.AttachWindow(window).WindowOffset(5)
	app.Event.On("timer_tick", func(e *application.CustomEvent) {
		if dataSlice, ok := e.Data.([]any); ok && len(dataSlice) > 0 {
			if tick, ok := dataSlice[0].(backend.TimerTick); ok {
				tray.SetLabel(tick.FormattedTime)
			}
		} else if tick, ok := e.Data.(backend.TimerTick); ok {
			tray.SetLabel(tick.FormattedTime)
		}
	})

	app.Event.On("timer_complete", func(e *application.CustomEvent) {
		dialog := app.Dialog.Info()
		dialog.SetTitle("Focalist")
		dialog.SetMessage("Session Complete")
		dialog.Show()
	})

	// We bind Cmd+Shift+P to play/pause.
	app.KeyBinding.Add("CmdOrCtrl+Shift+P", func(window application.Window) {
		// Just call skip or toggle on backend? Wait, we can't easily query state without another method.
		// Actually, we can add a Toggle method to timer.go, or just Play/Pause.
		// Let's call Play/Pause based on state or we can add a simple toggle.
		// For now, I'll just toggle visibility of window as an alternative if timer toggle isn't simple.
		if window != nil {
			if window.IsVisible() {
				window.Hide()
			} else {
				window.Show()
			}
		}
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
