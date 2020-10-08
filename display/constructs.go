package display

// type StatusScreen struct {
// 	Volume   string
// 	Tempo    string
// 	Title    string
// 	Progress float64
// 	State    string
// }

// func (l *Display) RenderStatus(s StatusScreen) {
// 	l.Clear()
// 	if s.State == "PLAYING" {
// 		l.SetColor(0, 255, 0)
// 		l.buffer[0] = "PLAYING:"
// 		l.buffer[1] = s.Title
// 		l.buffer[2] = l.progressBar(s.Progress)
// 		l.buffer[3] = "TEMPO: " + s.Tempo
// 	} else if s.State == "PAUSED" {
// 		l.SetColor(255, 60, 0)
// 		l.buffer[0] = "PLAYING:"
// 		l.buffer[1] = s.Title
// 		l.buffer[2] = l.progressBar(s.Progress)
// 		l.buffer[3] = "TEMPO: " + s.Tempo
// 	} else {
// 		l.buffer[1] = "     ATP MT-420"
// 		l.buffer[2] = " FLOPPY MIDI PLAYER"
// 	}
// 	l.render()
// }

// func (l *Display) RenderList(items []string, selector int) {

// 	l.Clear()

// 	page := int(math.Floor(float64(selector) / 4))
// 	pages := int(math.Floor(float64(len(items)) / 4))
// 	if len(items)%4 != 0 {
// 		pages++
// 	}

// 	list := items[page*4:]
// 	if len(list) > 4 {
// 		list = list[:4]
// 	}

// 	for i, item := range list {
// 		if i != selector%4 {
// 			l.buffer[i] = l.trim(item)
// 		} else {
// 			l.buffer[i] = l.trim("-> " + item)
// 		}
// 	}

// 	l.render()
// }
