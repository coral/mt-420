package main

import (
	"fmt"
	"io/ioutil"

	"github.com/baol/go-firmata"
	"github.com/coral/fluidsynth2"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting MT-420")

	// fl := floppy.New("/dev/fd0", "/media/floppy")
	// fl.Watch()

	s := fluidsynth2.NewSettings()
	synth := fluidsynth2.NewSynth(s)
	synth.SFLoad("files/SC-55.sf2", false)

	player := fluidsynth2.NewPlayer(synth)

	// dat, err := ioutil.ReadFile("/media/floppy/" + fl.FileIndex[0].Name())
	// if err != nil {
	// 	panic(err)
	// }

	dat, err := ioutil.ReadFile("./files/NEIL_YOUNG_-_Rockin_in_the_free_world.mid")
	if err != nil {
		panic(err)
	}
	player.AddMem(dat)
	player.Stop()
	// Easy way to set audio backend
	//s.SetString("audio.driver", "coreaudio")

	fluidsynth2.NewAudioDriver(s, synth)

	c, err := firmata.NewClient("/dev/cu.usbmodem14444301", 57600)
	if err != nil {
		panic("Cannot open client")
	}
	c.SetPinMode(2, firmata.Input)
	c.SetPinMode(3, firmata.Input)
	//fmt.Println(c.EnableAnalogInput(0, true))
	fmt.Println(c.EnableDigitalInput(2, true))
	fmt.Println(c.EnableDigitalInput(3, true))

	d := c.GetValues()

	for {
		m := <-d

		if player.GetStatus() == "PLAYING" {
			c.DigitalWrite(12, true)
		} else {
			c.DigitalWrite(12, false)
		}

		if m.IsAnalog() {
			// pin, val, err := m.GetAnalogValue()
			// if err != nil {
			// 	fmt.Println(err)
			// } else {
			// 	//fmt.Println(pin, val)
			// }
		} else {
			_, input, _ := m.GetDigitalValue()
			if input[2] == true {
				player.Stop()
				player.Seek(0)
				player.Play()
			}

			if input[3] == true {
				if player.GetStatus() == "PLAYING" {
					player.Stop()
					player.Seek(0)
				}
			}
		}
	}

	//player.Play()
	//player.Join()
}
