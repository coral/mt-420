package controller

import (
	"github.com/coral/mt-420/floppy"
	"github.com/coral/mt-420/lcd"
	"github.com/coral/mt-420/panel"
	"github.com/coral/mt-420/player"
)

var Options = map[string]Option{
	"status": &Status{},
}

type Controller struct {
	player  *player.Player
	panel   *panel.Panel
	floppy  *floppy.Disk
	display *lcd.LCD

	state string
}

func New(player *player.Player, panel *panel.Panel, floppy *floppy.Disk, l *lcd.LCD) *Controller {
	return &Controller{
		player:  player,
		panel:   panel,
		floppy:  floppy,
		display: l,
		state:   "status",
	}
}

func (c *Controller) Start() {

	quit := make(chan bool)
	evProp := make(chan string)

	go func() {
		msgChan := c.panel.GetEvents()
		for {
			e := <-msgChan
			if e == "escape" {
				quit <- true
			} else {
				evProp <- e
			}
		}
	}()

	for {
		Options[c.state].Run(c, evProp, quit)
	}

	//c.display.RenderList([]string{"passport.mid", "pepega.mid", "super-fighta.mid", "denis.mid"}, 2)
	// c.display.RenderStatus(
	// 	lcd.StatusScreen{
	// 		Title:    "passport.mid",
	// 		Tempo:    c.player.GetBPM(),
	// 		Volume:   "100%",
	// 		Progress: 50,
	// 	},
	// )

	// c.player.Play("files/passport.mid")
	// c.player.Wait()

}

func (c *Controller) handleEvents() {

}
