package player

import (
	"os"
	"strconv"

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
	SoundFont    string
	AudioBackend string
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
}

func New(c Configuration) (*Player, error) {
	sd := fluidsynth2.NewSettings()
	sd.SetString("audio.driver", c.AudioBackend)

	sy := fluidsynth2.NewSynth(sd)
	lf, err := sy.SFLoad(c.SoundFont, false)
	if err != nil {
		return nil, err
	}

	return &Player{
		Config:     c,
		Settings:   sd,
		synth:      sy,
		driver:     fluidsynth2.NewAudioDriver(sd, sy),
		fsPlayer:   fluidsynth2.NewPlayer(sy),
		loadedFont: lf,
	}, nil
}

func (p *Player) Play(filename string, data []byte) error {

	p.fsPlayer.Stop()
	if p.loaded {
		p.fsPlayer.Close()
		p.loaded = false
	}
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

	err := p.synth.SFUnload(p.loadedFont, true)
	if err != nil {
		return err
	}

	lf, err := p.synth.SFLoad("files/soundfonts/"+f.Name(), false)
	if err != nil {
		return err
	}
	p.loadedFont = lf
	p.state = "READY"

	return nil
}

func (p *Player) GetBPM() string {
	return strconv.Itoa(p.fsPlayer.GetBPM()) + " BPM"
}

func (p *Player) GetPlayingSong() string {
	pm := p.fsPlayer.GetStatus()
	if pm == "PLAYING" || pm == "DONE" {
		return p.filename
	}

	return ""

}

func (p *Player) GetState() string {
	st := p.fsPlayer.GetStatus()
	if st == "DONE" && p.state == PAUSED {
		return string(PAUSED)
	}
	p.state = State(st)
	return string(p.state)
}

func (p *Player) GetProgress() float64 {
	curr := p.fsPlayer.GetCurrentTick()
	total := p.fsPlayer.GetTotalTicks()
	if total > 0 {
		return float64(curr) / float64(total) * 100
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
