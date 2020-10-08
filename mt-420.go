package main

import (
	"flag"

	"github.com/coral/mt-420/config"
	"github.com/coral/mt-420/display"
	"github.com/coral/mt-420/display/terminal"
	"github.com/coral/mt-420/panel/kb"
	"github.com/coral/mt-420/panel/mtpanel"
	"github.com/coral/mt-420/storage"

	"github.com/coral/mt-420/controller"
	"github.com/coral/mt-420/display/lcd"
	"github.com/coral/mt-420/panel"
	"github.com/coral/mt-420/player"
	"github.com/coral/mt-420/storage/floppy"
	"github.com/coral/mt-420/storage/mock"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/errors/fmt"
)

func main() {

	///////////////////////////////////////////
	//FLAGS
	///////////////////////////////////////////
	virtual := flag.Bool("virtual", false, "virtual")
	terminalDisplay := flag.Bool("tdisp", false, "terminal display")
	mockFS := flag.Bool("mock", false, "mock")
	mockPath := flag.String("mockPath", "files/midi", "where to look for mock filesystem")
	configPath := flag.String("config", "config.json", "path to configuration")
	flag.Parse()

	///////////////////////////////////////////
	//LOAD CONFIG
	///////////////////////////////////////////
	lconfig, err := config.Load(*configPath)
	if err != nil {
		panic("could not load config : " + *configPath)
	}

	delay := lconfig.Lcd.BootDelay

	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)

	///////////////////////////////////////////
	//DISPLAY
	///////////////////////////////////////////
	var dp display.Device
	if *terminalDisplay {
		dp = terminal.New()
	} else {
		dp = lcd.New(lconfig.Lcd.Device)
	}
	err = dp.Init()
	if err != nil {
		panic(fmt.Sprintln("Display Fail", lconfig.Lcd.Device, err))
	}
	dp.SetContrast(lconfig.Lcd.Contrast)
	dp.SetBrightness(lconfig.Lcd.Brightness)

	display.DelayMessage(dp, "Starting MT-420", delay)

	///////////////////////////////////////////
	// Panel
	//////////////////////////////////////////
	display.DelayMessage(dp, "Connecting to panel", delay)
	var frontPanel panel.Panel
	if *virtual {
		frontPanel = kb.New()
		err = frontPanel.Init()
		if err != nil {
			panic(fmt.Sprintln("Panel Fail", lconfig.Panel.Device, err))
		}
	} else {
		frontPanel = mtpanel.New(lconfig.Panel.Device, lconfig.Panel.Baud, log)
		err = frontPanel.Init()
		if err != nil {
			panic(fmt.Sprintln("Panel Fail", lconfig.Panel.Device, err))
		}
		for _, be := range lconfig.Panel.Buttons {
			frontPanel.AddButton(be.Name, be.Message)
		}

		frontPanel.AddEncoder(panel.Encoder{
			Left:  byte(lconfig.Panel.Rotary.EncoderLeft),
			Right: byte(lconfig.Panel.Rotary.EncoderRight),
		})
	}

	///////////////////////////////////////////
	// Fluidsynth
	//////////////////////////////////////////
	display.DelayMessage(dp, "Loading Fluidsynth", delay)
	p, err := player.New(player.Configuration{
		SoundBank:        lconfig.Fluidsynth.Soundfonts,
		DefaultSoundFont: lconfig.Fluidsynth.Default,
		Gain:             lconfig.Fluidsynth.Gain,
	})
	defer p.Close()
	if err != nil {
		panic(err)
	}

	///////////////////////////////////////////
	// Floppy
	//////////////////////////////////////////
	display.DelayMessage(dp, "Warming up floppy", delay)
	var storage storage.Storage
	if *mockFS {
		storage = mock.New(*mockPath)
	} else {
		storage = floppy.New(lconfig.Floppy.Device, lconfig.Floppy.Mountpoint)
	}
	storage.Init()

	///////////////////////////////////////////
	// Controller
	//////////////////////////////////////////
	controller := controller.New(p, frontPanel, storage, dp)
	controller.Start()

}
