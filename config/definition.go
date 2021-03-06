package config

type C struct {
	Panel struct {
		Device  string   `json:"device"`
		Baud    int      `json:"baud"`
		Buttons []Button `json:"buttons"`
		Rotary  struct {
			EncoderLeft  int `json:"encoderLeft"`
			EncoderRight int `json:"encoderRight"`
		} `json:"rotary"`
	} `json:"panel"`
	Lcd struct {
		Device     string `json:"device"`
		Baud       int    `json:"baud"`
		BootDelay  int    `json:"bootdelay"`
		Contrast   int    `json:"contrast"`
		Brightness int    `json:"brightness"`
	} `json:"lcd"`
	Fluidsynth struct {
		Soundfonts string  `json:"soundfonts"`
		Default    string  `json:"default"`
		Gain       float32 `json:"gain"`
	} `json:"fluidsynth"`
	Floppy struct {
		Device     string `json:"device"`
		Mountpoint string `json:"mountpoint"`
	} `json:"floppy"`
}

type Button struct {
	Name    string `json:"name"`
	Message int    `json:"message,omitempty"`
}
