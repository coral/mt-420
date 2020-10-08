package terminal

import (
	"fmt"
	"os"
	"os/exec"
)

const eb = "                    "

type Terminal struct {
	buffer [4]string
}

func New() *Terminal {
	return &Terminal{
		buffer: [4]string{eb, eb, eb, eb},
	}
}

func (t *Terminal) Init() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}

func (t *Terminal) GetBuffer() [4]string {
	return t.buffer
}
func (t *Terminal) WriteBuffer(buf [4]string) {
	t.buffer = buf
	t.Render()
}

func (t *Terminal) ClearBuffer() {
	t.buffer = [4]string{eb, eb, eb, eb}
}

func (t *Terminal) WriteAt(x int, y int, str string) {
	bl := t.buffer[x-1][:y]

	bl = bl + str

	t.buffer[x-1] = bl
	t.Render()
}

func (t *Terminal) Render() {
	fmt.Printf("\033[0;0H")
	for _, val := range t.buffer {
		fmt.Println(val)
	}
}
func (t *Terminal) Clear() {
	t.buffer = [4]string{eb, eb, eb, eb}
	t.Render()
}

func (t *Terminal) SetColor(r byte, g byte, b byte) {

}
func (t *Terminal) SetContrast(c int) {

}
func (t *Terminal) SetBrightness(b int) {

}
