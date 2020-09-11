package panel

import (
	"github.com/baol/go-firmata"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/errors/fmt"
)

type Panel struct {
	cf     *firmata.FirmataClient
	device string
	layout Layout
}

type Layout struct {
	Buttons map[uint]Button
}

type Button struct {
	Name  string
	Pin   uint
	Value bool
}

func New(device string) *Panel {
	return &Panel{
		device: device,
		layout: Layout{
			Buttons: make(map[uint]Button),
		},
	}
}

func (p *Panel) Init() error {
	c, err := firmata.NewClient(p.device, 57600)
	if err != nil {
		return err
	}

	p.cf = c

	go p.scan(c.GetValues())

	return nil
}

func (p *Panel) AddButton(name string, pin int) {
	p.layout.Buttons[uint(pin)] = Button{
		Name: name,
		Pin:  uint(pin),
	}
	p.cf.SetPinMode(byte(pin), firmata.Input)
	p.cf.EnableDigitalInput(uint(pin), true)
}

func (p *Panel) scan(v <-chan firmata.FirmataValue) {
	for {
		update := <-v
		if update.IsAnalog() {

		} else {
			_, u, err := update.GetDigitalValue()
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("GetDigitalValue error")
			}

			for i, val := range u {
				if b, ok := p.layout.Buttons[uint(i)]; ok {
					if b.Value != val {
						b.Value = val.(bool)
						if val.(bool) {
							fmt.Println(b)
						}
					}
				}
			}
		}
	}
}
