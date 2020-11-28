package controller

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/coral/mt-420/display"
)

type SoundFonts struct {
}

func (m *SoundFonts) Run(c *Controller, events <-chan string, end chan bool) string {

	update := make(chan bool, 1)
	update <- true
	files := m.listDir(c.player.Config.SoundBank)

	selector := NewSelector(0, 1, 0, len(files)-1)

	var fn []string

	for _, file := range files {
		fn = append(fn, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
	}

	var renderEnd = make(chan bool)
	go func() {
		for {
			select {
			case <-renderEnd:
				return
			case <-update:
				display.RenderList(c.display, fn, selector.value)
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
			case "encoderRight":
				if selector.Increment() {
					update <- true
				}
			case "encoderLeft":
				if selector.Decrement() {
					update <- true
				}
			case "encoderClick":
				if len(files) > 0 {
					err := c.player.SwitchSoundFont(files[selector.Value()])
					if err != nil {
						display.Error(c.display, err)
					}
				}
				renderEnd <- true
				return "status"
			}
		}
	}
}

func (m *SoundFonts) Name() string {
	return "soundfonts"
}

func (m *SoundFonts) listDir(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var t []os.FileInfo

	for _, file := range files {
		if !file.IsDir() {
			if filepath.Ext(file.Name()) == ".sf2" {
				t = append(t, file)
			}
		}
	}

	return t
}
