package kb

import (
	"os"

	"github.com/coral/mt-420/panel"
	"github.com/eiannone/keyboard"
)

type KB struct {
	events chan string
}

func New() *KB {
	return &KB{
		events: make(chan string, 1),
	}
}

func (p *KB) Init() error {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}

	go p.key()

	return nil
}

func (p *KB) GetEvents() <-chan string {
	return p.events
}

func (p *KB) key() {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch r := string(char); r {
		case "d":
			p.events <- "encoderRight"
		case "a":
			p.events <- "encoderLeft"
		case "w":
			p.events <- "encoderClick"
		case "p":
			p.events <- "play"
		case "o":
			p.events <- "pause"
		case "i":
			p.events <- "stop"
		case "u":
			p.events <- "menu"
		case "y":
			p.events <- "escape"
		}
		if key == keyboard.KeyEsc {
			os.Exit(1)
		}
	}
}

func (p *KB) AddButton(name string, pin int) {}

func (p *KB) AddEncoder(e panel.Encoder) {}

func (p *KB) SetEncoderColor(r int, g int, b int) {}
