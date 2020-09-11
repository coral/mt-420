package main

import (
	"github.com/coral/mt-420/panel"
	"github.com/coral/mt-420/player"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/errors/fmt"
)

func main() {
	log.Info("Starting MT-420")

	panel := panel.New("/dev/cu.usbmodem14444301")
	err := panel.Init()
	if err != nil {
		panic(err)
	}

	panel.AddButton("play", 2)
	panel.AddButton("stop", 3)

	events := panel.GetEvents()
	go func() {
		for {
			e := <-events
			fmt.Println(e)
		}
	}()

	p, err := player.New(player.Configuration{
		SoundFont:    "files/SC-55.sf2",
		AudioBackend: "coreaudio",
	})
	defer p.Close()
	if err != nil {
		panic(err)
	}

	err = p.Play("files/NEIL_YOUNG_-_Rockin_in_the_free_world.mid")
	if err != nil {
		panic(err)
	}

	p.Wait()
}
