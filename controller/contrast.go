package controller

import "github.com/coral/mt-420/display"

type Contrast struct {
}

func (m *Contrast) Run(c *Controller, events <-chan string, end chan bool) string {

	min := 100
	max := 255

	update := make(chan bool, 1)
	update <- true

	selector := NewSelector(c.display.GetContrast(), 5, min, max)

	var renderEnd = make(chan bool)
	go func() {
		for {
			select {
			case <-renderEnd:
				return
			case <-update:
				display.RenderSlider(c.display, "Set Contrast", selector.Value(), min, max)
			}
		}
	}()

	for {
		select {
		case <-end:
			renderEnd <- true
			return "status"
		case m := <-events:
			switch m {
			case "encoderRight":
				if selector.Increment() {
					c.display.SetContrast(selector.Value())
					update <- true
				}
			case "encoderLeft":
				if selector.Decrement() {
					c.display.SetContrast(selector.Value())
					update <- true
				}
			case "encoderClick":
				renderEnd <- true
				return "settings"
			}
		}
	}
}

func (m *Contrast) Name() string {
	return "volume"
}
