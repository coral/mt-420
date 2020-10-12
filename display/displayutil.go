package display

import "math"

func GetEmptyBuffer() [4]string {
	const eb = "                    "
	return [4]string{eb, eb, eb, eb}
}

func CenterText(s string) string {
	const eb = "                    "
	a := 10 - math.Floor(float64(len(s)/2)) - 1
	return eb[:int(a)] + s
}
