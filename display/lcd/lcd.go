package lcd

import (
	"strings"
	"time"

	"github.com/tarm/serial"
)

const eb = "                    "

type LCD struct {
	device       string
	contrast     int
	brightness   int
	buffer       [4]string
	serialConfig *serial.Config
	conn         *serial.Port
	lastrender   time.Time
}

func New(device string) *LCD {
	return &LCD{
		device:     device,
		contrast:   255,
		brightness: 255,
		buffer:     [4]string{eb, eb, eb, eb},
	}
}

func (l *LCD) Init() error {

	//This code inits the USB + Serial RGB Backlight Character LCD Backpack
	//from Adafruit.
	// https://learn.adafruit.com/usb-plus-serial-backpack/command-reference

	l.serialConfig = &serial.Config{Name: l.device, Baud: 115200}
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
	s.Write([]byte{0xFE, 0x99, byte(l.brightness)})
	time.Sleep(10 * time.Millisecond)

	//Set contrast
	s.Write([]byte{0xFE, 0x50, byte(l.contrast)})
	time.Sleep(10 * time.Millisecond)

	////Set backlight
	s.Write([]byte{0xFE, 0xD0, 0xFF, 0xFF, 0x00})
	time.Sleep(10 * time.Millisecond)

	return nil
}

///////////////////////////////////////////
//BUFFER FUNCTIONS
///////////////////////////////////////////

func (l *LCD) GetBuffer() [4]string {
	return l.buffer
}

func (l *LCD) WriteBuffer(newBuf [4]string) {

	// This function is mainly because the performance of the LCD backpack from
	// Adafruit is absymal. For that reason I'm trying to be smart about how I
	// choose to render the different parts.

	type diff struct {
		r    byte
		line int
		pos  int
	}

	//Trim buffer
	for i, m := range newBuf {
		//fmt.Println(i, m, len(l.trim(m)))
		newBuf[i] = l.trim(m)
	}

	//DIFF buffer
	var changes []diff
	var changedLines []int

	for i, m := range newBuf {
		if m != l.buffer[i] {
			changedLines = append(changedLines, i)
			for si, bs := range l.buffer[i] {
				if m[si] != byte(bs) {
					changes = append(changes, diff{
						r:    m[si],
						line: i,
						pos:  si,
					})
				}
			}
		}
	}

	//Short changelist, writing individual chars
	if len(changes) < 10 {
		for _, c := range changes {
			l.WriteAt(c.line+1, c.pos, string(c.r))
		}
		l.buffer = newBuf
		return
	}

	//Medium change, writing full screen
	if len(changedLines) < 3 {

		for _, line := range changedLines {
			l.WriteAt(line+1, 0, string(newBuf[line]))
		}
		l.buffer = newBuf
		return
	}

	//Big change, full clear
	l.buffer = newBuf
	l.Clear()
	time.Sleep(5 * time.Millisecond)
	l.Render()

}

func (l *LCD) ClearBuffer() {
	l.buffer = [4]string{eb, eb, eb, eb}
}

///////////////////////////////////////////
// WRITING
///////////////////////////////////////////

func (l *LCD) WriteAt(x int, y int, str string) {
	bl := l.buffer[x-1][:y]
	bl = bl + str
	l.buffer[x-1] = bl
	l.Render()
}

func (l *LCD) Render() {

	//Jump to start
	l.conn.Write([]byte{0xFE, 0x48})
	time.Sleep(10 * time.Millisecond)
	for _, val := range l.buffer {
		l.conn.Write([]byte(l.trim(val)))
		time.Sleep(6 * time.Millisecond)
	}

	l.lastrender = time.Now()

}

func (l *LCD) Clear() {
	l.conn.Write([]byte{0xFE, 0x58})
}

///////////////////////////////////////////
// SETTINGS
///////////////////////////////////////////

func (l *LCD) SetColor(r byte, g byte, b byte) {
	////Set backlight
	l.conn.Write([]byte{0xFE, 0xD0, r, g, b})
	time.Sleep(10 * time.Millisecond)
}

func (l *LCD) SetContrast(c int) {
	l.contrast = c
	//Set contrast
	l.conn.Write([]byte{0xFE, 0x50, byte(l.contrast)})
	time.Sleep(10 * time.Millisecond)
}

func (l *LCD) SetBrightness(b int) {

	l.brightness = b
	//Set contrast
	l.conn.Write([]byte{0xFE, 0x99, byte(l.brightness)})
	time.Sleep(10 * time.Millisecond)
}

///////////////////////////////////////////
// INTERNAL
///////////////////////////////////////////

func (l *LCD) writeLine(line string, index int) {
	l.buffer[index] = line
}

func (l *LCD) trim(si string) string {
	//If longer than 20, trim
	if len(si) > 20 {
		return si[0:20]
	}
	//If shorter than 20, extend
	if len(si) < 20 {
		var b strings.Builder
		b.WriteString(si)
		for i := 0; i < 20-len(si); i++ {
			b.WriteString(" ")
		}

		return b.String()
	}

	//If 20, we gucchimucchi
	return si
}
