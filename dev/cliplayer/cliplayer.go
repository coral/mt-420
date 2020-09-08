package main

import (
	"flag"
	"io/ioutil"
	"runtime"

	"github.com/coral/fluidsynth2"
)

func main() {

	soundfont := flag.String("soundfont", "/Users/coral/Music/Soundfonts/sc55.sf2", "soundfont file")
	midifile := flag.String("midi", "", "midi file")

	flag.Parse()

	s := fluidsynth2.NewSettings()
	if runtime.GOOS == "darwin" {
		s.SetString("audio.driver", "coreaudio")
	} else {
		s.SetString("audio.driver", "alsa")
	}
	synth := fluidsynth2.NewSynth(s)
	synth.SFLoad(*soundfont, false)

	player := fluidsynth2.NewPlayer(synth)

	dat, err := ioutil.ReadFile(*midifile)
	if err != nil {
		panic(err)
	}

	player.AddMem(dat)
	fluidsynth2.NewAudioDriver(s, synth)

	player.Play()
	player.Join()
}
