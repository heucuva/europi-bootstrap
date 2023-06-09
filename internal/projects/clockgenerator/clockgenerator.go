package main

import (
	"time"

	europi "github.com/awonak/EuroPiGo"
	bootstrap "github.com/heucuva/europi-bootstrap"
	"github.com/heucuva/europi-bootstrap/experimental/screenbank"
	"github.com/heucuva/europi-bootstrap/internal/projects/clockgenerator/module"
	"github.com/heucuva/europi-bootstrap/internal/projects/clockgenerator/screen"
)

type application struct {
	clock *module.ClockGenerator

	ui             *screenbank.ScreenBank
	screenMain     screen.Main
	screenSettings screen.Settings
}

func newApplication() (*application, error) {
	clock := &module.ClockGenerator{}
	app := &application{
		clock: clock,

		screenMain: screen.Main{
			Clock: clock,
		},
		screenSettings: screen.Settings{
			Clock: clock,
		},
	}

	var err error
	app.ui, err = screenbank.NewScreenBank(
		screenbank.WithScreen("main", "\u2b50", &app.screenMain),
		screenbank.WithScreen("settings", "\u2611", &app.screenSettings),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *application) Start(e *europi.EuroPi) {
	if err := app.clock.Init(module.Config{
		BPM:          120.0,
		GateDuration: time.Millisecond * 100,
		Enabled:      true,
		ClockOut: func(value bool) {
			if value {
				e.CV1.SetCV(1.0)
			} else {
				e.CV1.SetCV(0.0)
			}
			bootstrap.ForceRepaintUI(e)
		},
	}); err != nil {
		panic(err)
	}
}

func (app *application) MainLoop(e *europi.EuroPi, deltaTime time.Duration) {
	app.clock.Tick(deltaTime)
}

func main() {
	app, err := newApplication()
	if err != nil {
		panic(err)
	}

	pi := europi.New()

	// some options shown below are being explicitly set to their defaults
	// only to showcase their existence.
	if err := bootstrap.Bootstrap(
		pi,
		bootstrap.EnableDisplayLogger(false),
		bootstrap.InitRandom(true),
		bootstrap.App(
			app,
			bootstrap.AppMainLoopInterval(time.Millisecond*1),
		),
		bootstrap.UI(
			app.ui,
			bootstrap.UIRefreshRate(time.Millisecond*50),
		),
	); err != nil {
		panic(err)
	}
}
