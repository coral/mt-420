package controller

import (
	"time"

	"github.com/coral/mt-420/lcd"
)

type Status struct {
}

func (m *Status) Run(c *Controller, events <-chan string, end chan bool) string {

	var renderEnd = make(chan bool)
	go func() {
		for {
			select {
			case <-renderEnd:
				return
			default:
				c.display.RenderStatus(
					lcd.StatusScreen{
						Title:    c.player.GetPlayingSong(),
						Tempo:    c.player.GetBPM(),
						Volume:   "100%",
						Progress: c.player.GetProgress(),
					},
				)
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
			case "play":
				c.player.Play("files/passport.mid")
			}
		}
	}
}

func (m *Status) Name() string {
	return "status"
}
