package main

import (
	"flag"
	"io"
	"os"
	"runtime"
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
	mockFS := flag.Bool("mock", false, "mockfs")
	flag.Parse()
	delay := 0

	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)

	//LCD
	display := lcd.New(true, log)
	delayWriter("Starting MT-420", delay, display)

	ersh := ErrorShim{
		lcd: display,
	}
	mw := io.MultiWriter(os.Stdout, ersh)
	log.SetOutput(mw)

	delayWriter("Connecting to panel", delay, display)
	// Panel
	panel := panel.New("/dev/cu.usbmodem14444301", *virtual, log)
	err := panel.Init()
	if err != nil {
		panic(err)
	}

	panel.AddButton("play", 2)
	panel.AddButton("escape", 3)

	delayWriter("Loading Fluidsynth", delay, display)

	var backend string

	os := runtime.GOOS
	switch os {
	case "darwin":
		backend = "coreaudio"
	case "linux":
		backend = "alsa"
	}

	//Player
	p, err := player.New(player.Configuration{
		SoundFont:    "files/SC-55.sf2",
		AudioBackend: backend,
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
	controller := controller.New(p, panel, storage, display)
	controller.Start()

}

func delayWriter(message string, delay int, l *lcd.LCD) {
	m := time.Duration(delay)
	l.Message(message)
	time.Sleep(m * time.Second)
}
