package lcd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
)

type LCD struct {
	logger  *logrus.Logger
	device  string
	virtual bool
	buffer  [4]string
}

type StatusScreen struct {
	Volume   string
	Tempo    string
	Title    string
	Progress int
	State    string
}

func New(virtual bool, log *logrus.Logger) *LCD {
	if virtual {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	return &LCD{
		logger:  log,
		virtual: virtual,
		buffer:  [4]string{"", "", "", ""},
	}
}

//MENUS

func (l *LCD) RenderStatus(s StatusScreen) {
	l.Clear()
	l.buffer[0] = "NOW PLAYING:"
	l.buffer[1] = s.Title
	l.buffer[2] = strconv.Itoa(s.Progress)
	l.buffer[3] = "TEMPO: " + s.Tempo + " VOL:" + s.Volume
	l.render()
}

func (l *LCD) RenderList(items []string, selector int) {
	for i, item := range items {
		if i != selector {
			l.buffer[i] = item
		} else {
			l.buffer[i] = "-> " + item
		}
	}

	l.render()
}

///FUNCTIONS

func (l *LCD) WriteBuffer(buf [4]string) {
	l.buffer = buf
}

func (l *LCD) Update() {
	l.render()
}

func (l *LCD) Error(err error) {
	errMess := err.Error()
	l.Clear()
	l.writeLine("ERROR", 0)
	l.writeLine(errMess, 1)

}

func (l *LCD) Message(message string) {
	l.Clear()
	l.logger.Info(message)
	l.writeLine(message, 0)
	l.render()
}

func (l *LCD) Clear() {
	for i, _ := range l.buffer {
		l.buffer[i] = "                    "
	}
	l.render()
}

func (l *LCD) writeLine(line string, index int) {
	l.buffer[index] = line
}

func (l *LCD) render() {

	if l.virtual {
		fmt.Printf("\033[0;0H")
		for _, val := range l.buffer {
			fmt.Println(val)
		}
	}
}
