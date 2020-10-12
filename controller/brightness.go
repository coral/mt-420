package controller

import "github.com/coral/mt-420/display"

type Brightness struct {
}

func (m *Brightness) Run(c *Controller, events <-chan string, end chan bool) string {

	min := 5
	max := 255

	update := make(chan bool, 1)
	update <- true

	selector := NewSelector(c.display.GetBrightness(), 5, min, max)

	var renderEnd = make(chan bool)
	go func() {
		for {
			select {
			case <-renderEnd:
				return
			case <-update:
				display.RenderSlider(c.display, "Set Brightness", selector.Value(), min, max)
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
					c.display.SetBrightness(selector.Value())
					update <- true
				}
			case "encoderLeft":
				if selector.Decrement() {
					c.display.SetBrightness(selector.Value())
					update <- true
				}
			case "encoderClick":
				renderEnd <- true
				return "settings"
			}
		}
	}
}

func (m *Brightness) Name() string {
	return "volume"
}
