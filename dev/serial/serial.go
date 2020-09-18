package main

import (
	"fmt"
	"time"

	firmata "github.com/coral/go-firmata"
)

func main() {
	c, err := firmata.NewClient("/dev/cu.usbmodem14444301", 57600)
	if err != nil {
		panic("Cannot open client")
	}

	fmt.Println("HELLO")
	//	c.SetPinMode(2, firmata.Input)
	//c.SetPinMode(3, firmata.Input)
	fmt.Println(c.SetPinMode(6, firmata.PWM))

	err = c.EnableEncoder(1, 2, 3)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(c.EnableAnalogInput(12, true))
	fmt.Println(c.EnableDigitalInput(12, true))
	//fmt.Println(c.EnableDigitalInput(3, true))
	var i = 0
	for {
		c.AnalogWrite(6, byte(i))
		i++
		if i == 254 {
			i = 0
		}
		//c.DigitalWrite(11, true)
		time.Sleep(10 * time.Millisecond)
	}

	// d := c.GetValues()

	// for {
	// 	m := <-d

	// 	if m.IsAnalog() {
	// 		// pin, val, err := m.GetAnalogValue()
	// 		// if err != nil {
	// 		// 	fmt.Println(err)
	// 		// } else {
	// 		// 	//fmt.Println(pin, val)
	// 		// }
	// 	} else {
	// 		_, input, _ := m.GetDigitalValue()
	// 		if input[2] == true {
	// 			c.DigitalWrite(12, true)
	// 		} else {
	// 			c.DigitalWrite(12, false)
	// 		}
	// 	}
	// }

}
