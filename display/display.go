package display

type Device interface {
	Init() error

	GetBuffer() [4]string
	WriteBuffer(buf [4]string)
	ClearBuffer()

	WriteAt(x int, y int, str string)
	Render()
	Clear()

	SetColor(r byte, g byte, b byte)
	SetContrast(c int)
	SetBrightness(b int)
}

func Message(d Device, message string) {

	d.Clear()
	d.ClearBuffer()
	d.WriteAt(1, 1, message)
}

func Error(d Device, err error) {
	d.SetColor(255, 0, 0)
	d.Clear()
	d.ClearBuffer()
	d.WriteAt(1, 1, "     - ERROR -")
	d.WriteAt(2, 1, err.Error())
}
