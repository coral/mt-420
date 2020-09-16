package main

import (
	"fmt"
	"math"
)

func main() {

	denis := []float64{0.00, 3.2, 10.12, 33.4, 40, 60, 70, 88, 99.9, 100}

	for _, d := range denis {
		fmt.Println(progressBar(d), len(progressBar(d)))
	}
}

func progressBar(x float64) string {
	var ic float64 = 0.18
	var pg int
	pg = int(math.Ceil(ic * x))
	if x > 95 {
		pg = int(math.Floor(ic * x))
	}
	template := []rune("I                  I")
	for i := 0; i < pg; i++ {

		template[i+1] = '-'
		if pg == 1 {
			template[i+1] = '>'
		} else if pg != 18 {
			template[i+2] = '>'
		}

	}

	return string(template)
}
