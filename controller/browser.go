package controller

import (
	"path/filepath"
	"strings"
	"time"
)

type Browser struct {
}

func (m *Browser) Run(c *Controller, events <-chan string, end chan bool) string {

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
			default:
				if len(fn) > 0 {
					c.display.RenderList(fn, selector.value)
				} else {
					c.display.Message("No files on floppy")
				}
				time.Sleep(50 * time.Millisecond)
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
				selector.Increment()
			case "encoderLeft":
				selector.Decrement()
			case "encoderClick":
				if len(files) > 0 {
					d, err := c.storage.LoadFile(files[selector.Value()])
					if err != nil {
						c.display.Error(err)
					}
					c.player.Play(files[selector.Value()].Name(), d)
				}
				renderEnd <- true
				return "status"
			}
		}
	}
}

func (m *Browser) Name() string {
	return "browser"
}
