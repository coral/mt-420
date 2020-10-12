package controller

import (
	"time"

	"github.com/coral/mt-420/display"
)

type Settings struct {
}

func (m *Settings) Run(c *Controller, events <-chan string, end chan bool) string {

	c.display.SetColor(0, 128, 255)
	c.panel.SetEncoderColor(0, 120, 255)
	update := make(chan bool, 1)
	update <- true

	set := []string{"Volume", "Soundfonts", "Brightness", "Contrast", "Return to player"}

	selector := NewSelector(0, 1, 0, len(set)-1)

	var renderEnd = make(chan bool)
	go func() {
		for {
			select {
			case <-renderEnd:
				return
			case <-update:
				display.RenderList(c.display, set, selector.value)
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
					update <- true
				}
			case "encoderLeft":
				if selector.Decrement() {
					update <- true
				}
			case "encoderClick":
				rt := "status"
				switch set[selector.Value()] {
				case "Volume":
					rt = "volume"
				case "Soundfonts":
					rt = "soundfonts"
				case "Brightness":
					rt = "brightness"
				case "Contrast":
					rt = "contrast"
				case "Return to player":
					rt = "status"
				}
				renderEnd <- true
				time.Sleep(20 * time.Millisecond)
				return rt
			}
		}
	}
}

func (m *Settings) Name() string {
	return "settings"
}
