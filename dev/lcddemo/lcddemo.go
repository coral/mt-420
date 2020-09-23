package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/serial1", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	//s.Write([]byte{0xFE, 0xD1, 0x20, 0x04})

	s.Write([]byte{0xFE, 0x58})
	time.Sleep(10 * time.Millisecond)
	s.Write([]byte{0xFE, 0x54})
	time.Sleep(10 * time.Millisecond)
	s.Write([]byte{0xFE, 0x99, 0xFF})
	time.Sleep(10 * time.Millisecond)
	s.Write([]byte{0xFE, 0x50, 220})

	time.Sleep(10 * time.Millisecond)

	s.Write([]byte{0xFE, 0xD0, 0xFF, 0xFF, 0x00})

	// var str string
	// str = str + "8HlKzqD6aREu6Ptwv59G"
	// str = str + "4V5uXxfQPWj909pmDTdv"
	// str = str + "tq3Ky5X0A4U8xGIrUlQV"
	// str = str + "jM4oIX9KQPUvYXdqpCrU"

	s.Write([]byte("8HlKzqD6aREu"))
	time.Sleep(20 * time.Millisecond)
	s.Write([]byte{0x0D, 0x0A})
	time.Sleep(20 * time.Millisecond)
	s.Write([]byte("4V5uXxfQPWj909pmDTdv"))
	time.Sleep(10 * time.Millisecond)
	//s.Write([]byte("tq3Ky5X0A4U8xGIrUlQV"))
	time.Sleep(10 * time.Millisecond)
	//s.Write([]byte("jM4oIX9KQPUvYXdqpCrU"))

	// for {
	// 	time.Sleep(500 * time.Millisecond)
	// 	s.Write([]byte{0xFE, 0xD0, 0x00, 0xFF, 0xFF})
	// 	time.Sleep(500 * time.Millisecond)
	// 	s.Write([]byte{0xFE, 0xD0, 0xFF, 0xFF, 0xFF})
	// 	time.Sleep(500 * time.Millisecond)
	// 	s.Write([]byte{0xFE, 0xD0, 0x00, 0x00, 0xFF})
	// 	time.Sleep(500 * time.Millisecond)
	// 	s.Write([]byte{0xFE, 0xD0, 0xFF, 0x00, 0x00})
	// 	time.Sleep(500 * time.Millisecond)
	// 	s.Write([]byte{0xFE, 0xD0, 0xFF, 0xFF, 0x00})
	// }

	buffer = [4]string{"", "", "", ""}
	Clear()
	buffer[1] = "     ATP MT-420     "
	buffer[2] = " FLOPPY MIDI PLAYER "

	s.Write([]byte{0xFE, 0x40})
	time.Sleep(50 * time.Millisecond)
	for _, dm := range buffer {
		time.Sleep(50 * time.Millisecond)
		s.Write([]byte(dm))
		fmt.Println(dm)
	}

}

var buffer [4]string

func Clear() {
	for i, _ := range buffer {
		buffer[i] = "                    "
	}
	//l.render()
}

func render() {

	//Clear LCD
	var bu []byte
	for _, val := range buffer {
		bu = append(bu, []byte(trim(val))...)

	}

	fmt.Println(bu)

}

func trim(si string) string {
	if len(si) > 20 {
		return si[0:19]
	}
	return si
}
