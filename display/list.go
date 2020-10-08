package display

import "math"

func RenderList(d Device, items []string, selector int) {

	buf := GetEmptyBuffer()

	page := int(math.Floor(float64(selector) / 4))
	pages := int(math.Floor(float64(len(items)) / 4))
	if len(items)%4 != 0 {
		pages++
	}

	list := items[page*4:]
	if len(list) > 4 {
		list = list[:4]
	}

	for i, item := range list {
		if i != selector%4 {
			buf[i] = item
		} else {
			buf[i] = "-> " + item
		}
	}

	d.WriteBuffer(buf)
}
