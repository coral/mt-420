package panel

import (
	"os"

	"github.com/eiannone/keyboard"
	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type Panel struct {
	config  *serial.Config
	port    *serial.Port
	logger  *logrus.Logger
	device  string
	layout  Layout
	events  chan string
	virtual bool
}

type Layout struct {
	Buttons map[byte]Button
	Enc     Encoder
}

type Button struct {
	Name  string
	Value byte
}

type Encoder struct {
	Name  string
	Left  byte
	Right byte
	LED   RGBLED
}

type RGBLED struct {
	Red   byte
	Green byte
	Blue  byte
}

func New(device string, virtual bool, logger *logrus.Logger) *Panel {

	return &Panel{
		device: device,
		logger: logger,
		layout: Layout{
			Buttons: make(map[byte]Button),
		},
		events:  make(chan string, 1),
		virtual: virtual,
	}
}

func (p *Panel) Init() error {
	if !p.virtual {
		p.config = &serial.Config{Name: p.device, Baud: 115200}
		port, err := serial.OpenPort(p.config)
		if err != nil {
			return err
		}

		p.port = port

		go p.scan()
	} else {
		if err := keyboard.Open(); err != nil {
			panic(err)
		}

		go p.key()
	}
	return nil
}

func (p *Panel) AddButton(name string, pin int) {
	p.layout.Buttons[byte(pin)] = Button{
		Name:  name,
		Value: byte(pin),
	}
}

func (p *Panel) AddEncoder(e Encoder) {
	p.layout.Enc = e
}

func (p *Panel) SetEncoderColor(r int, g int, b int) {

	if p.layout.Enc.LED.Red != byte(clamp(r)) || p.layout.Enc.LED.Green != byte(clamp(g)) || p.layout.Enc.LED.Blue != byte(clamp(b)) {
		p.layout.Enc.LED.Red = byte(clamp(r))
		p.layout.Enc.LED.Green = byte(clamp(g))
		p.layout.Enc.LED.Blue = byte(clamp(b))

		p.writeEncoderRGB()
	}
}

func (p *Panel) writeEncoderRGB() {
	p.port.Write([]byte{0xCA, 0xFE, 0xBA, 0xBE, p.layout.Enc.LED.Red, p.layout.Enc.LED.Green, p.layout.Enc.LED.Blue})
}

func (p *Panel) scan() {

	buf := make([]byte, 1)
	for {
		n, _ := p.port.Read(buf)
		if n > 0 {

			if b, ok := p.layout.Buttons[buf[0]]; ok {
				p.events <- b.Name
			}

			if p.layout.Enc.Left == buf[0] {
				p.events <- "encoderLeft"
			}
			if p.layout.Enc.Right == buf[0] {
				p.events <- "encoderRight"
			}
		}
	}
}

func (p *Panel) key() {
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

func (p *Panel) GetEvents() <-chan string {
	return p.events
}

func clamp(i int) int {
	if i < 0 {
		return 0
	}
	if i > 255 {
		return 255
	}
	return i
}
