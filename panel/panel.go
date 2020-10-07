package panel

type Layout struct {
	Buttons map[byte]Button
	Enc     Encoder
}

type Button struct {
	Name  string
	Value byte
}

type Encoder struct {
	Name  string
	Left  byte
	Right byte
	LED   RGBLED
}

type RGBLED struct {
	Red   byte
	Green byte
	Blue  byte
}

type Panel interface {
	AddButton(name string, pin int)
	AddEncoder(e Encoder)
	SetEncoderColor(r int, g int, b int)

	Init() error
	GetEvents() <-chan string
}
