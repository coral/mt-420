package controller

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/coral/mt-420/display"
)

type Browser struct {
}

func (m *Browser) Run(c *Controller, events <-chan string, end chan bool) string {

	c.display.SetColor(0, 255, 255)
	c.panel.SetEncoderColor(0, 120, 255)
	update := make(chan bool, 1)
	update <- true
	files := c.storage.ListFiles()

	selector := NewSelector(0, 1, 0, len(files)-1)

	var fn []string

	for _, file := range files {
		fn = append(fn, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
	}

	var renderEnd = make(chan bool)
	go func() {
		for {
			select {
			case <-renderEnd:
				return
			case <-update:
				if len(fn) > 0 {
					display.RenderList(c.display, fn, selector.value)
				} else {
					c.display.SetColor(255, 128, 0)
					display.Message(c.display, "No files on floppy")
				}
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
				if len(files) > 0 {
					d, err := c.storage.LoadFile(files[selector.Value()])
					if err != nil {
						display.Error(c.display, err)
					}
					err = c.player.Play(files[selector.Value()].Name(), d)
					if err != nil {
						display.Error(c.display, err)
					}
				}
				renderEnd <- true
				time.Sleep(50 * time.Millisecond)
				return "status"
			}
		}
	}
}

func (m *Browser) Name() string {
	return "browser"
}
