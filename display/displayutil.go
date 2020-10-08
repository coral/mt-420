package display

func GetEmptyBuffer() [4]string {
	const eb = "                    "
	return [4]string{eb, eb, eb, eb}
}
