package player

import (
	"os"
	"path"
	"runtime"

	"github.com/coral/fluidsynth2"
)

type State string

const (
	PLAYING State = "PLAYING"
	READY   State = "READY"
	PAUSED  State = "PAUSED"
	DONE    State = "DONE"
	UNKNOWN State = "UNKNOWN"
)

type Configuration struct {
	SoundBank        string
	DefaultSoundFont string
	AudioBackend     string
	Gain             float32
}

type Player struct {
	Config     Configuration
	Settings   fluidsynth2.Settings
	synth      fluidsynth2.Synth
	driver     fluidsynth2.AudioDriver
	fsPlayer   fluidsynth2.Player
	state      State
	filename   string
	loaded     bool
	loadedFont int
	fontname   string
}

func New(c Configuration) (*Player, error) {
	sd := fluidsynth2.NewSettings()
	if c.AudioBackend == "" {
		os := runtime.GOOS
		switch os {
		case "darwin":
			c.AudioBackend = "coreaudio"
		case "linux":
			c.AudioBackend = "alsa"
		}
	}
	sd.SetString("audio.driver", c.AudioBackend)

	sy := fluidsynth2.NewSynth(sd)
	lf, err := sy.SFLoad(path.Join(c.SoundBank, c.DefaultSoundFont), false)
	if err != nil {
		return nil, err
	}

	sy.SetGain(c.Gain)

	return &Player{
		Config:     c,
		Settings:   sd,
		synth:      sy,
		driver:     fluidsynth2.NewAudioDriver(sd, sy),
		fsPlayer:   fluidsynth2.NewPlayer(sy),
		loadedFont: lf,
		fontname:   c.DefaultSoundFont,
	}, nil
}

func (p *Player) Play(filename string, data []byte) error {

	p.Stop()
	p.fsPlayer = fluidsynth2.NewPlayer(p.synth)

	p.filename = filename

	errCode := p.fsPlayer.AddMem(data)
	if errCode != nil {
		return errCode
	}
	p.loaded = true

	return nil

}

func (p *Player) Pause() {
	pm := p.fsPlayer.GetStatus()
	if pm == "PLAYING" {
		p.fsPlayer.Stop()
		p.state = PAUSED
	} else if pm == "DONE" && p.state == PAUSED {
		p.fsPlayer.Play()
	}

}

func (p *Player) Stop() {

	p.fsPlayer.Stop()
	if p.loaded {
		p.fsPlayer.Close()
		p.loaded = false
	}
	p.state = State(p.fsPlayer.GetStatus())

}

func (p *Player) SwitchSoundFont(f os.FileInfo) error {

	//RESET SHIT
	p.fsPlayer.Stop()
	if p.loaded {
		p.fsPlayer.Close()
		p.loaded = false
	}

	if path.Base(f.Name()) != p.fontname {

		err := p.synth.SFUnload(p.loadedFont, true)
		if err != nil {
			return err
		}

		lf, err := p.synth.SFLoad(path.Join(p.Config.SoundBank, f.Name()), false)
		if err != nil {
			return err
		}

		p.loadedFont = lf
	}
	p.state = "READY"

	return nil
}

func (p *Player) GetBPM() int {
	if p.loaded {
		return p.fsPlayer.GetBPM()
	}
	return 0
}

func (p *Player) ChangeTempo(t int) {
	v := p.fsPlayer.GetBPM()
	p.fsPlayer.SetBPM(v + t)
}

func (p *Player) GetPlayingSong() string {
	if p.loaded {
		pm := p.fsPlayer.GetStatus()
		if pm == "PLAYING" || pm == "DONE" {
			return p.filename
		}
	}

	return ""

}

func (p *Player) GetState() string {
	if p.loaded {
		st := p.fsPlayer.GetStatus()
		if st == "DONE" && p.state == PAUSED {
			return string(PAUSED)
		}
		p.state = State(st)
		return string(p.state)
	}
	return string(DONE)
}

func (p *Player) GetProgress() float64 {
	if p.loaded {
		curr := p.fsPlayer.GetCurrentTick()
		total := p.fsPlayer.GetTotalTicks()
		if total > 0 {
			return float64(curr) / float64(total) * 100
		}
	}
	return 0
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
