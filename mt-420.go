package main

import (
	"time"

	"github.com/coral/mt-420/player"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting MT-420")

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

	time.Sleep(6 * time.Second)

	err = p.Play("files/passport.mid")
	if err != nil {
		panic(err)
	}

	p.Wait()
}
