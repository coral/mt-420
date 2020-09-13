package panel

import (
	"fmt"
	"os"

	"github.com/baol/go-firmata"
	"github.com/eiannone/keyboard"
	"github.com/sirupsen/logrus"
)

type Panel struct {
	cf      *firmata.FirmataClient
	logger  *logrus.Logger
	device  string
	layout  Layout
	events  chan string
	virtual bool
}

type Layout struct {
	Buttons map[uint]Button
	LEDs    map[string]LED
	RGB     map[string]RGBLED
}

type Button struct {
	Name  string
	Pin   uint
	Value bool
}

type LED struct {
	Name string
	Pin  uint
}

type RGBLED struct {
	Name     string
	RedPin   uint
	GreenPin uint
	BluePin  uint
}

func New(device string, virtual bool, logger *logrus.Logger) *Panel {

	return &Panel{
		device: device,
		logger: logger,
		layout: Layout{
			Buttons: make(map[uint]Button),
			LEDs:    make(map[string]LED),
		},
		events:  make(chan string, 1),
		virtual: virtual,
	}
}

func (p *Panel) Init() error {
	if !p.virtual {
		c, err := firmata.NewClient(p.device, 57600)
		if err != nil {
			return err
		}

		p.cf = c

		go p.scan(c.GetValues())
	} else {
		if err := keyboard.Open(); err != nil {
			panic(err)
		}

		go p.key()
	}
	return nil
}

func (p *Panel) AddButton(name string, pin int) {
	p.layout.Buttons[uint(pin)] = Button{
		Name: name,
		Pin:  uint(pin),
	}
	if !p.virtual {
		p.cf.SetPinMode(byte(pin), firmata.Input)
		p.cf.EnableDigitalInput(uint(pin), true)
	}
}

func (p *Panel) AddLED(name string, pin int) {
	p.layout.LEDs[name] = LED{
		Name: name,
		Pin:  uint(pin),
	}
	p.cf.SetPinMode(byte(pin), firmata.PWM)
}

func (p *Panel) SetLEDIntensity(name string, intensity byte) error {
	if v, ok := p.layout.LEDs[name]; ok {
		p.cf.AnalogWrite(v.Pin, intensity)
	}
	return fmt.Errorf("No LED defined with this name")
}

func (p *Panel) AddRGBLED(name string, rpin int, gpin int, bpin int) {
	p.layout.RGB[name] = RGBLED{
		Name:     name,
		RedPin:   uint(rpin),
		GreenPin: uint(gpin),
		BluePin:  uint(bpin),
	}
	p.cf.SetPinMode(byte(rpin), firmata.PWM)
	p.cf.SetPinMode(byte(gpin), firmata.PWM)
	p.cf.SetPinMode(byte(bpin), firmata.PWM)
}

func (p *Panel) SetRGBIntensity(name string, r byte, g byte, b byte) error {
	if v, ok := p.layout.RGB[name]; ok {
		p.cf.AnalogWrite(v.RedPin, r)
		p.cf.AnalogWrite(v.GreenPin, g)
		p.cf.AnalogWrite(v.BluePin, b)
	}
	return fmt.Errorf("No RGB LED defined with this name")
}

func (p *Panel) scan(v <-chan firmata.FirmataValue) {
	for {
		update := <-v
		if update.IsAnalog() {

		} else {
			_, u, err := update.GetDigitalValue()
			if err != nil {
				p.logger.WithFields(logrus.Fields{
					"Error": err,
				}).Error("GetDigitalValue error")
			}

			for i, val := range u {
				if b, ok := p.layout.Buttons[uint(i)]; ok {
					if b.Value != val {
						b.Value = val.(bool)
						if val.(bool) {
							p.events <- b.Name
						}
					}
				}
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
		fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)

		switch r := string(char); r {
		case "d":
			p.events <- "el"
		case "a":
			p.events <- "er"
		case "p":
			p.events <- "play"
		case "q":
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
