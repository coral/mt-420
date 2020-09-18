package main

import (
	"flag"
	"io"
	"os"
	"time"

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
	virtual := flag.Bool("virtual", false, "virtual")
	terminalDisplay := flag.Bool("tdisp", false, "terminal display")
	mockFS := flag.Bool("mock", false, "mock")

	flag.Parse()
	delay := 0

	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)

	//LCD
	display := lcd.New("/dev/cu.usbmodem142444301", *terminalDisplay, log)
	err := display.Init()
	if err != nil {
		panic("display fail")
	}
	delayWriter("Starting MT-420", delay, display)

	ersh := ErrorShim{
		lcd: display,
	}
	mw := io.MultiWriter(os.Stdout, ersh)
	log.SetOutput(mw)

	delayWriter("Connecting to panel", delay, display)
	// Panel
	frontPanel := panel.New("/dev/cu.usbmodem14444301", *virtual, log)
	err = frontPanel.Init()
	if err != nil {
		panic(err)
	}

	frontPanel.AddButton("play", 0x20)
	frontPanel.AddButton("pause", 0x21)
	frontPanel.AddButton("stop", 0x22)
	frontPanel.AddButton("menu", 0x23)
	frontPanel.AddButton("escape", 0x24)
	frontPanel.AddButton("encoderClick", 0x25)

	frontPanel.AddEncoder(panel.Encoder{
		Left:  0x26,
		Right: 0x27,
	})
	frontPanel.SetEncoderColor(0, 255, 0)

	delayWriter("Loading Fluidsynth", delay, display)

	//Player
	p, err := player.New(player.Configuration{
		SoundBank:        "files/soundfonts",
		DefaultSoundFont: "Roland SC-55.sf2",
	})
	defer p.Close()
	if err != nil {
		panic(err)
	}

	delayWriter("Warming up floppy", delay, display)
	//Floppy
	var storage storage.Storage
	if *mockFS {
		storage = mock.New("files/midi")
	} else {
		storage = floppy.New("/dev/sdb", "/media/floppy")
	}
	storage.Init()

	//Controller
	controller := controller.New(p, frontPanel, storage, display)
	controller.Start()

}

func delayWriter(message string, delay int, l *lcd.LCD) {
	m := time.Duration(delay)
	l.Message(message)
	time.Sleep(m * time.Second)
}
