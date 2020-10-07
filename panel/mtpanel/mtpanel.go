package mtpanel

import (
	"github.com/coral/mt-420/panel"
	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type Panel struct {
	config *serial.Config
	port   *serial.Port
	logger *logrus.Logger
	device string
	layout panel.Layout
	events chan string
	baud   int
}

func New(device string, baud int, logger *logrus.Logger) *Panel {

	return &Panel{
		device: device,
		baud:   baud,
		logger: logger,
		layout: panel.Layout{
			Buttons: make(map[byte]panel.Button),
		},
		events: make(chan string, 1),
	}
}

func (p *Panel) Init() error {
	p.config = &serial.Config{Name: p.device, Baud: p.baud}
	port, err := serial.OpenPort(p.config)
	if err != nil {
		return err
	}

	p.port = port

	go p.scan()

	return nil
}

func (p *Panel) AddButton(name string, pin int) {
	p.layout.Buttons[byte(pin)] = panel.Button{
		Name:  name,
		Value: byte(pin),
	}
}

func (p *Panel) AddEncoder(e panel.Encoder) {
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
