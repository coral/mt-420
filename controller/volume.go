package controller

import "github.com/coral/mt-420/display"

type Volume struct {
}

func (m *Volume) Run(c *Controller, events <-chan string, end chan bool) string {

	min := 0.0
	max := 3.0

	update := make(chan bool, 1)
	update <- true

	selector := NewFloatSelector(float64(c.player.GetGain()), 0.1, min, max)

	var renderEnd = make(chan bool)
	go func() {
		for {
			select {
			case <-renderEnd:
				return
			case <-update:
				display.RenderFloatSlider(c.display, "Volume", selector.Value(), min, max)
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
					c.player.SetGain(float32(selector.Value()))
					update <- true
				}
			case "encoderLeft":
				if selector.Decrement() {
					c.player.SetGain(float32(selector.Value()))
					update <- true
				}
			case "encoderClick":

				renderEnd <- true
				return "settings"
			}
		}
	}
}

func (m *Volume) Name() string {
	return "volume"
}
