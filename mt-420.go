package main

import (
	"io"
	"os"

	"github.com/coral/mt-420/controller"
	"github.com/coral/mt-420/floppy"
	"github.com/coral/mt-420/lcd"
	"github.com/coral/mt-420/panel"
	"github.com/coral/mt-420/player"
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
	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)

	//LCD
	display := lcd.New(true, log)

	ersh := ErrorShim{
		lcd: display,
	}
	mw := io.MultiWriter(os.Stdout, ersh)
	log.SetOutput(mw)

	display.Message("Starting MT-420")

	// Panel
	panel := panel.New("/dev/cu.usbmodem14444301", log)
	err := panel.Init()
	if err != nil {
		panic(err)
	}

	panel.AddButton("play", 2)
	panel.AddButton("stop", 3)

	// events := panel.GetEvents()
	// go func() {
	// 	for {
	// 		e := <-events
	// 		fmt.Println(e)
	// 	}
	// }()

	//Player
	p, err := player.New(player.Configuration{
		SoundFont:    "files/SC-55.sf2",
		AudioBackend: "coreaudio",
	})
	defer p.Close()
	if err != nil {
		panic(err)
	}

	//Floppy
	fl := floppy.New("/dev/fd0", "/media/floppy")

	//Controller
	controller := controller.New(p, panel, fl, display)
	controller.Start()

}
