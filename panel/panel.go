package panel

import (
	"github.com/baol/go-firmata"
	"github.com/sirupsen/logrus"
)

type Panel struct {
	cf     *firmata.FirmataClient
	logger *logrus.Logger
	device string
	layout Layout
	events chan string
}

type Layout struct {
	Buttons map[uint]Button
}

type Button struct {
	Name  string
	Pin   uint
	Value bool
}

func New(device string, logger *logrus.Logger) *Panel {

	return &Panel{
		device: device,
		logger: logger,
		layout: Layout{
			Buttons: make(map[uint]Button),
		},
		events: make(chan string, 1),
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

func (p *Panel) GetEvents() <-chan string {
	return p.events
}
