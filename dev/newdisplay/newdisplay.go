package main

import (
	"fmt"
	"time"

	"github.com/coral/mt-420/display"
	"github.com/coral/mt-420/display/lcd"
)

//this is just to test the refactor of the display lib.

const eb = "                    "

func main() {

	eb := [4]string{eb, eb, eb, eb}

	lcd := lcd.New("/dev/cu.usbmodem142444401")

	//lcd := terminal.New()

	err := lcd.Init()
	if err != nil {
		panic(err)
	}

	lcd.SetContrast(200)

	time.Sleep(1 * time.Second)

	b := lcd.GetBuffer()

	b[0] = "   BENIS sssssssssssssssssssssssssssssssss   "

	lcd.WriteBuffer(b)

	time.Sleep(2 * time.Second)

	lcd.WriteBuffer(b)

	time.Sleep(2 * time.Second)

	eb[0] = "HALKKKASFOJKF"

	lcd.WriteBuffer(eb)

	time.Sleep(2 * time.Second)

	lcd.WriteAt(2, 3, "DENISSSSS")

	time.Sleep(2 * time.Second)

	display.Message(lcd, "HELLO BOYS")

	progress := 0
	for {
		display.RenderStatus(lcd, display.StatusScreen{
			Tempo:    "120",
			Title:    "KKONA BOYS",
			Progress: float64(progress),
			State:    "PLAYING",
		})
		time.Sleep(50 * time.Millisecond)
		if progress > 99 {
			break
		}
		progress++
	}

	item := 0
	for {

		display.RenderList(lcd, []string{"hej man",
			"tjena boys halloo",
			"jaoooooooooooo",
			"cperik"}, item)

		if item == 3 {
			break
		}
		item++
		time.Sleep(800 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)

	display.RenderSlider(lcd, "INT", 25, 0, 100)

	time.Sleep(2 * time.Second)

	display.RenderFloatSlider(lcd, "FLOAT", 26.7, 0, 87.3)

	time.Sleep(2 * time.Second)

	display.Error(lcd, fmt.Errorf("KORV KORV"))

}
