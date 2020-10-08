package controller

import (
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coral/mt-420/display"
)

type Status struct {
	c *Controller
}

func (m *Status) Run(c *Controller, events <-chan string, end chan bool) string {
	m.c = c
	var renderEnd = make(chan bool)

	if c.player.GetState() == "DONE" {
		m.render()
	}

	go func() {
		var progress = 0
		var state string = c.player.GetState()
		for {
			select {
			case <-renderEnd:
				return
			default:
				var ic float64 = 0.18
				pg := int(math.Ceil(ic * c.player.GetProgress()))
				st := c.player.GetState()

				switch st {
				case "PLAYING":
					c.panel.SetEncoderColor(0, 255, 0)
				case "PAUSED":
					c.panel.SetEncoderColor(255, 255, 0)
				case "DONE":
					c.panel.SetEncoderColor(255, 255, 255)
				case "READY":
					c.panel.SetEncoderColor(255, 0, 0)
				}

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
		case ev := <-events:
			switch ev {
			case "play":
				//c.player.Play("files/passport.mid")
			case "encoderClick":
				renderEnd <- true
				return "browser"
			case "encoderRight":
				c.player.ChangeTempo(+5)
				time.Sleep(10 * time.Millisecond)
				m.render()
			case "encoderLeft":
				c.player.ChangeTempo(-5)
				time.Sleep(10 * time.Millisecond)
				m.render()
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
	t := strconv.Itoa(m.c.player.GetBPM())
	t = t + " BPM"
	go display.RenderStatus(m.c.display,
		display.StatusScreen{
			Title: strings.TrimSuffix(
				m.c.player.GetPlayingSong(),
				filepath.Ext(m.c.player.GetPlayingSong())),
			Tempo:    t,
			Progress: m.c.player.GetProgress(),
			State:    m.c.player.GetState(),
		},
	)
}

// func (m *Status) renderTempo() {
// 	t := strconv.Itoa(m.c.player.GetBPM())
// 	t = t + " BPM  "
// 	m.c.display.WriteFrom(4, 8, t)
// }
