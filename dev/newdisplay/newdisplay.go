package main

import (
	"time"

	"github.com/coral/mt-420/display/terminal"
)

//this is just to test the refactor of the display lib.

const eb = "                    "

func main() {

	eb := [4]string{eb, eb, eb, eb}

	//l := logrus.New()

	//lcd := lcd.New("/dev/cu.usbmodem142444401")

	lcd := terminal.New()

	err := lcd.Init()
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	b := lcd.GetBuffer()

	b[2] = "   BENIS s   "

	lcd.WriteBuffer(b)

	time.Sleep(1 * time.Second)

	lcd.WriteBuffer(b)

	time.Sleep(1 * time.Second)

	eb[0] = "HALKKKASFOJKF"

	lcd.WriteBuffer(eb)

	time.Sleep(1 * time.Second)

	lcd.WriteAt(2, 3, "DENISSSSS")
}
