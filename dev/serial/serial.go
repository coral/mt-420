package main

import (
	"fmt"

	firmata "github.com/baol/go-firmata"
)

func main() {
	c, err := firmata.NewClient("/dev/cu.usbmodem14444301", 57600)
	if err != nil {
		panic("Cannot open client")
		panic(err)
	}
	c.SetPinMode(2, firmata.Input)
	c.SetPinMode(3, firmata.Input)
	//fmt.Println(c.EnableAnalogInput(0, true))
	fmt.Println(c.EnableDigitalInput(2, true))
	fmt.Println(c.EnableDigitalInput(3, true))

	d := c.GetValues()

	for {
		m := <-d

		if m.IsAnalog() {
			// pin, val, err := m.GetAnalogValue()
			// if err != nil {
			// 	fmt.Println(err)
			// } else {
			// 	//fmt.Println(pin, val)
			// }
		} else {
			_, input, _ := m.GetDigitalValue()
			if input[2] == true {
				c.DigitalWrite(12, true)
			} else {
				c.DigitalWrite(12, false)
			}
		}
	}

}
