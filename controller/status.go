package controller

import (
	"math"
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
	go c.display.RenderStatus(
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
	go func() {
		var progress = 0
		for {
			select {
			case <-renderEnd:
				return
			default:
				var ic float64 = 0.18
				pg := int(math.Ceil(ic * c.player.GetProgress()))
				if progress != pg {
					progress = pg
					go c.display.RenderStatus(
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
				} else {
					time.Sleep(100 * time.Millisecond)
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
