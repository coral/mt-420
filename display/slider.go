package display

import "strconv"

func RenderFloatSlider(d Device, title string, v float64, min float64, max float64) {
	buf := GetEmptyBuffer()

	buf[0] = CenterText(title)
	buf[2] = progressBar(v / 0.03)
	buf[3] = CenterText(strconv.FormatFloat(v, 'g', 2, 64))

	d.WriteBuffer(buf)

}
