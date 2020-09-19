package main

import (
	"flag"
	"io"
	"os"
	"time"

	"github.com/coral/mt-420/config"
	"github.com/coral/mt-420/storage"

	"github.com/coral/mt-420/controller"
	"github.com/coral/mt-420/lcd"
	"github.com/coral/mt-420/panel"
	"github.com/coral/mt-420/player"
	"github.com/coral/mt-420/storage/floppy"
	"github.com/coral/mt-420/storage/mock"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/errors/fmt"
)

type ErrorShim struct {
	lcd *lcd.LCD
}

func (e ErrorShim) Write(data []byte) (n int, err error) {
	e.lcd.Error(fmt.Errorf(string(data)))
	fmt.Println(data)
	return len(data), nil
}

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
	//LCD
	///////////////////////////////////////////
	display := lcd.New(lconfig.Lcd.Device, *terminalDisplay, log)
	err = display.Init()
	if err != nil {
		panic(fmt.Sprintln("Display Fail", lconfig.Lcd.Device, err))
	}
	delayWriter("Starting MT-420", delay, display)

	ersh := ErrorShim{
		lcd: display,
	}
	mw := io.MultiWriter(os.Stdout, ersh)
	log.SetOutput(mw)

	///////////////////////////////////////////
	// Panel
	//////////////////////////////////////////
	delayWriter("Connecting to panel", delay, display)
	frontPanel := panel.New(lconfig.Panel.Device, lconfig.Panel.Baud, *virtual, log)
	err = frontPanel.Init()
	if err != nil {
		panic(err)
	}
	for _, be := range lconfig.Panel.Buttons {
		frontPanel.AddButton(be.Name, be.Message)
	}

	frontPanel.AddEncoder(panel.Encoder{
		Left:  byte(lconfig.Panel.Rotary.EncoderLeft),
		Right: byte(lconfig.Panel.Rotary.EncoderRight),
	})

	///////////////////////////////////////////
	// Fluidsynth
	//////////////////////////////////////////
	delayWriter("Loading Fluidsynth", delay, display)
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
	delayWriter("Warming up floppy", delay, display)
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
	controller := controller.New(p, frontPanel, storage, display)
	controller.Start()

}

func delayWriter(message string, delay int, l *lcd.LCD) {
	m := time.Duration(delay)
	l.Message(message)
	time.Sleep(m * time.Second)
}
