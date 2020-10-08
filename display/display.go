package display

type Device interface {
	Init() error

	GetBuffer() [4]string
	WriteBuffer(buf [4]string)

	WriteAt(x int, y int, str string)
	Render()
	Clear()

	SetColor(r byte, g byte, b byte)
	SetContrast(c int)
	SetBrightness(b int)
}
