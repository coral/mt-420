package player

import (
	"fmt"
	"io/ioutil"

	"github.com/coral/fluidsynth2"
)

type Configuration struct {
	SoundFont    string
	AudioBackend string
}

type Player struct {
	Config   Configuration
	Settings fluidsynth2.Settings
	synth    fluidsynth2.Synth
	driver   fluidsynth2.AudioDriver
	fsPlayer fluidsynth2.Player
}

func New(c Configuration) (*Player, error) {
	sd := fluidsynth2.NewSettings()
	sd.SetString("audio.driver", c.AudioBackend)

	sy := fluidsynth2.NewSynth(sd)
	_, err := sy.SFLoad(c.SoundFont, false)
	if err != nil {
		return nil, err
	}

	return &Player{
		Config:   c,
		Settings: sd,
		synth:    sy,
		driver:   fluidsynth2.NewAudioDriver(sd, sy),
		fsPlayer: fluidsynth2.NewPlayer(sy),
	}, nil
}

func (p *Player) Play(filename string) error {

	p.fsPlayer.Stop()
	p.fsPlayer.Close()
	p.fsPlayer = fluidsynth2.NewPlayer(p.synth)

	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	errCode := p.fsPlayer.AddMem(dat)
	if errCode != nil {
		return fmt.Errorf("Could not load MIDI data into FluidSynth")
	}

	return nil

}

func (p *Player) Wait() {
	p.fsPlayer.Join()
}

func (p *Player) Close() {
	p.fsPlayer.Close()
	p.driver.Close()
	p.synth.Close()
	p.Settings.Close()
}
