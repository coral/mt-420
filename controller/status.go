package controller

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/coral/mt-420/lcd"
)

type Status struct {
}

func (m *Status) Run(c *Controller, events <-chan string, end chan bool) string {
	c.display.SetColor(0, 255, 0)
	var renderEnd = make(chan bool)
	go func() {
		for {
			select {
			case <-renderEnd:
				return
			default:
				c.display.RenderStatus(
					lcd.StatusScreen{
						Title: strings.TrimSuffix(
							c.player.GetPlayingSong(),
							filepath.Ext(c.player.GetPlayingSong())),
						Tempo:    c.player.GetBPM(),
						Volume:   "100%",
						Progress: c.player.GetProgress(),
						State:    c.player.GetState(),
					},
				)
				time.Sleep(1000 * time.Millisecond)
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
			case "play":
				//c.player.Play("files/passport.mid")
			case "encoderClick":
				renderEnd <- true
				return "browser"
			}
		}
	}
}

func (m *Status) Name() string {
	return "status"
}
