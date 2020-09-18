package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/cu.usbmodem14444301", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	go func() {

		// scanner := bufio.NewScanner(s)
		// for {
		// 	for scanner.Scan() {
		// 		fmt.Println(scanner.Text())
		// 	}
		// }
		buf := make([]byte, 1)
		for {
			n, err := s.Read(buf)
			if n > 0 {
				fmt.Println(buf[0])
			}
			if err == io.EOF {
				break
			}
		}
	}()

	time.Sleep(4 * time.Second)
	fmt.Println([]byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00})

	n, err := s.Write([]byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 30, 128})
	if err != nil {
		panic(err)
	}
	time.Sleep(50 * time.Millisecond)
	fmt.Println(n)

	m := make(chan bool)
	m <- true
}
