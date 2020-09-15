package lcd

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type LCD struct {
	logger       *logrus.Logger
	device       string
	virtual      bool
	buffer       [4]string
	serialConfig *serial.Config
	conn         *serial.Port
}

type StatusScreen struct {
	Volume   string
	Tempo    string
	Title    string
	Progress float64
	State    string
}

func New(device string, virtual bool, log *logrus.Logger) *LCD {
	return &LCD{
		logger:  log,
		device:  device,
		virtual: virtual,
		buffer:  [4]string{"", "", "", ""},
	}
}

func (l *LCD) Init() error {
	if l.virtual {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		l.serialConfig = &serial.Config{Name: "/dev/cu.usbmodem142444301", Baud: 9600}
		s, err := serial.OpenPort(l.serialConfig)
		if err != nil {
			return err
		}

		l.conn = s

		//Clear screen
		s.Write([]byte{0xFE, 0x58})
		time.Sleep(10 * time.Millisecond)

		//Block cursor off
		s.Write([]byte{0xFE, 0x54})
		time.Sleep(10 * time.Millisecond)

		//Autoscroll off
		s.Write([]byte{0xFE, 0x52})
		time.Sleep(10 * time.Millisecond)

		//Set Brightness
		s.Write([]byte{0xFE, 0x99, 0xFF})
		time.Sleep(10 * time.Millisecond)

		//Set contrast
		s.Write([]byte{0xFE, 0x50, 200})
		time.Sleep(10 * time.Millisecond)

		////Set backlight
		s.Write([]byte{0xFE, 0xD0, 0xFF, 0xFF, 0x00})
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

//MENUS

func (l *LCD) RenderStatus(s StatusScreen) {
	l.Clear()
	if s.State == "PLAYING" {
		l.buffer[0] = "PLAYING:"
		l.buffer[1] = s.Title
		l.buffer[2] = l.progressBar(s.Progress)
		l.buffer[3] = "TEMPO: " + s.Tempo
	} else {
		l.buffer[0] = "ATP MT-420"
	}
	l.render()
}

func (l *LCD) RenderList(items []string, selector int) {

	l.Clear()

	page := int(math.Floor(float64(selector) / 4))
	pages := int(math.Floor(float64(len(items)) / 4))
	if len(items)%4 != 0 {
		pages++
	}

	list := items[page*4:]
	if len(list) > 4 {
		list = list[:4]
	}

	for i, item := range list {
		if i != selector%4 {
			l.buffer[i] = l.trim(item)
		} else {
			l.buffer[i] = l.trim("-> " + item)
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

func (l *LCD) SetColor(r byte, g byte, b byte) {
	////Set backlight
	l.conn.Write([]byte{0xFE, 0xD0, r, g, b})
	time.Sleep(10 * time.Millisecond)
}

//////INTERNAL

func (l *LCD) writeLine(line string, index int) {
	l.buffer[index] = line
}

func (l *LCD) render() {

	if l.virtual {
		fmt.Printf("\033[0;0H")
		for _, val := range l.buffer {
			fmt.Println(val)
		}
	} else {
		//Clear LCD
		l.conn.Write([]byte{0xFE, 0x48})
		time.Sleep(10 * time.Millisecond)
		for _, val := range l.buffer {
			data := l.trim(val)
			if len(data) > 19 {
				l.conn.Write([]byte(l.trim(val)))
			} else {
				l.conn.Write(append([]byte(data), []byte{0x0D, 0x0A}...))
			}
		}
	}

}

func (l *LCD) trim(si string) string {
	if len(si) > 20 {
		return si[0:19]
	}
	return si
}

func (l *LCD) progressBar(x float64) string {
	var ic float64 = 0.18
	pg := int(math.Ceil(ic * x))
	template := "I                 I"
	for i := 0; i < pg; i++ {
		template = template[:i] + "-" + template[1:]
	}
	return template
}
