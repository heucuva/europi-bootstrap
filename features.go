package bootstrap

import (
	"log"
	"os"

	europi "github.com/awonak/EuroPiGo"
	"github.com/awonak/EuroPiGo/experimental/displaylogger"
	"github.com/awonak/EuroPiGo/hardware/hal"
)

var (
	dispLog displaylogger.Logger
)

func enableDisplayLogger(e europi.Hardware) {
	if dispLog != nil {
		// already enabled - can happen when panicking
		return
	}

	display := europi.Display(e)
	if display == nil {
		// no display, can't continue
		return
	}

	log.SetFlags(0)
	dispLog = displaylogger.NewLogger(display)
	log.SetOutput(dispLog)
}

func disableDisplayLogger(e europi.Hardware) {
	flushDisplayLogger(e)
	dispLog = nil
	log.SetOutput(os.Stdout)
}

func flushDisplayLogger(e europi.Hardware) {
	if dispLog != nil {
		dispLog.Flush()
	}
}

func initRandom(e europi.Hardware) {
	if rnd := e.Random(); rnd != nil {
		_ = rnd.Configure(hal.RandomGeneratorConfig{})
	}
}

func uninitRandom(e europi.Hardware) {
}
