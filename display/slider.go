package display

import "strconv"

func RenderFloatSlider(d Device, title string, v float64, min float64, max float64) {
	buf := GetEmptyBuffer()

	buf[0] = CenterText(title)
	buf[2] = progressBar(v / ((max - min) / 100))
	buf[3] = CenterText(strconv.FormatFloat(v, 'g', 2, 64))

	d.WriteBuffer(buf)

}
