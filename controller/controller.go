package controller

import (
	"github.com/coral/mt-420/display"
	"github.com/coral/mt-420/panel"
	"github.com/coral/mt-420/player"
	"github.com/coral/mt-420/storage"
)

var Options = map[string]Option{
	"status":     &Status{},
	"browser":    &Browser{},
	"soundfonts": &SoundFonts{},
	"settings":   &Settings{},
	"volume":     &Volume{},
}

type Controller struct {
	player  *player.Player
	panel   panel.Panel
	storage storage.Storage
	display display.Device

	state string
}

func New(player *player.Player, panel panel.Panel, storage storage.Storage, l display.Device) *Controller {
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
