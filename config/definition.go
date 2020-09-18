package config
a
type C struct {
	Panel struct {
		Buttons []struct {
			Name string `json:"name"`
			Gpio string `json:"gpio"`
		} `json:"buttons"`
		Rotary struct {
			Name string `json:"name"`
			Led  struct {
				Name  string `json:"name"`
				Red   string `json:"red"`
				Green string `json:"green"`
				Blue  string `json:"blue"`
			} `json:"led"`
			Button struct {
				Name string `json:"name"`
				Gpio string `json:"gpio"`
			} `json:"button"`
			EncoderLeft  string `json:"encoderLeft"`
			EncoderRight string `json:"encoderRight"`
		} `json:"rotary"`
	} `json:"panel"`
	Lcd struct {
		Device string `json:"device"`
		Baud   int    `json:"baud"`
	} `json:"lcd"`
	Fluidsynth struct {
		Soundfonts string `json:"soundfonts"`
		Default    string `json:"default"`
	} `json:"fluidsynth"`
	Floppy struct {
		Device     string `json:"device"`
		Mountpoint string `json:"mountpoint"`
	} `json:"floppy"`
}
