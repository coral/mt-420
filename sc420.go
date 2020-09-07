package main

import (
	"io/ioutil"

	"github.com/coral/fluidsynth2"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting SC-420")

	s := fluidsynth2.NewSettings()
	synth := fluidsynth2.NewSynth(s)
	i := synth.SFLoad("files/SC-55.sf2", false)

	player := fluidsynth2.NewPlayer(synth)
	//fmt.Println(player.Add("files/be_sharp_bw_redfarn.mid"))
	dat, err := ioutil.ReadFile("files/be_sharp_bw_redfarn.mid")
	if err != nil {
		panic(err)
	}

	player.AddMem(dat)

	s.SetString("audio.driver", "coreaudio")

	adriver := fluidsynth2.NewAudioDriver(s, synth)
	_ = adriver

	player.Play()
	player.Join()

}
