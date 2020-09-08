package main

import (
	"io/ioutil"

	"github.com/coral/fluidsynth2"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting MT-420")

	s := fluidsynth2.NewSettings()
	synth := fluidsynth2.NewSynth(s)
	synth.SFLoad("files/SC-55.sf2", false)

	player := fluidsynth2.NewPlayer(synth)

	dat, err := ioutil.ReadFile("files/NEIL_YOUNG_-_Rockin_in_the_free_world.mid")
	if err != nil {
		panic(err)
	}

	player.AddMem(dat)

	// Easy way to set audio backend
	//s.SetString("audio.driver", "coreaudio")

	fluidsynth2.NewAudioDriver(s, synth)

	player.Play()
	player.Join()
}
