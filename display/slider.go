package display

import "strconv"

func RenderSlider(d Device, title string, v int, min int, max int) {
	buf := GetEmptyBuffer()

	buf[0] = CenterText(title)
	buf[2] = progressBar(float64(v-min) / ((float64(max) - float64(min)) / 100))
	buf[3] = CenterText(strconv.Itoa(v))

	d.WriteBuffer(buf)

}

func RenderFloatSlider(d Device, title string, v float64, min float64, max float64) {
	buf := GetEmptyBuffer()

	buf[0] = CenterText(title)
	buf[2] = progressBar((v - min) / ((max - min) / 100))
	buf[3] = CenterText(strconv.FormatFloat(v, 'g', 2, 64))

	d.WriteBuffer(buf)

}
