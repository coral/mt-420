package controller

import (
	"github.com/coral/mt-420/lcd"
	"github.com/coral/mt-420/panel"
	"github.com/coral/mt-420/player"
	"github.com/coral/mt-420/storage"
)

var Options = map[string]Option{
	"status":     &Status{},
	"browser":    &Browser{},
	"soundfonts": &SoundFonts{},
}

type Controller struct {
	player  *player.Player
	panel   *panel.Panel
	storage storage.Storage
	display *lcd.LCD

	state string
}

func New(player *player.Player, panel *panel.Panel, storage storage.Storage, l *lcd.LCD) *Controller {
	return &Controller{
		player:  player,
		panel:   panel,
		storage: storage,
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
			switch e {
			case "escape":
				quit <- true
			case "pause":
				c.player.Pause()
			case "stop":
				c.player.Stop()
			default:
				evProp <- e
			}
		}
	}()

	for {
		c.state = Options[c.state].Run(c, evProp, quit)
	}

}

func (c *Controller) handleEvents() {

}
