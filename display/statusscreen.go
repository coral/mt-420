package display

import (
	"math"
)

type StatusScreen struct {
	Tempo    string
	Title    string
	Progress float64
	State    string
}

func RenderStatus(d Device, s StatusScreen) {
	buf := GetEmptyBuffer()
	if s.State == "PLAYING" {
		d.SetColor(0, 255, 0)
		buf[0] = "PLAYING:"
		buf[1] = s.Title
		buf[2] = progressBar(s.Progress)
		buf[3] = "TEMPO: " + s.Tempo
	} else if s.State == "PAUSED" {
		d.SetColor(255, 60, 0)
		buf[0] = "PLAYING:"
		buf[1] = s.Title
		buf[2] = progressBar(s.Progress)
		buf[3] = "TEMPO: " + s.Tempo
	} else {
		buf[1] = "     ATP MT-420"
		buf[2] = " FLOPPY MIDI PLAYER"
	}
	d.WriteBuffer(buf)

}

func progressBar(x float64) string {
	var ic float64 = 0.18
	var pg int
	pg = int(math.Ceil(ic * x))
	if x > 95 {
		pg = int(math.Floor(ic * x))
	}
	template := []rune("I                  I")
	for i := 0; i < pg; i++ {

		template[i+1] = '-'
		if pg == 1 {
			template[i+1] = '>'
		} else if pg != 18 {
			template[i+2] = '>'
		}

	}

	return string(template)
}
