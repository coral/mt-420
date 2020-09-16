package controller

import (
	"math"
	"path/filepath"
	"strings"
	"time"

	"github.com/coral/mt-420/lcd"
)

type Status struct {
	c *Controller
}

func (m *Status) Run(c *Controller, events <-chan string, end chan bool) string {
	//c.display.SetColor(0, 255, 0)
	m.c = c
	var renderEnd = make(chan bool)
	m.render()
	go func() {
		var progress = 0
		var state string = "DONE"
		for {
			select {
			case <-renderEnd:
				return
			default:
				var ic float64 = 0.18
				pg := int(math.Ceil(ic * c.player.GetProgress()))
				st := c.player.GetState()

				if progress != pg && state != "PAUSED" {
					progress = pg
					m.render()
				}

				if state != st {
					state = st
					m.render()
				}

				time.Sleep(100 * time.Millisecond)
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
			case "menu":
				renderEnd <- true
				return "soundfonts"
			}
		}
	}
}

func (m *Status) Name() string {
	return "status"
}

func (m *Status) render() {
	go m.c.display.RenderStatus(
		lcd.StatusScreen{
			Title: strings.TrimSuffix(
				m.c.player.GetPlayingSong(),
				filepath.Ext(m.c.player.GetPlayingSong())),
			Tempo:    m.c.player.GetBPM(),
			Volume:   "100%",
			Progress: m.c.player.GetProgress(),
			State:    m.c.player.GetState(),
		},
	)
}
