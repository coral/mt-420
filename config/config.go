package config

import (
	"encoding/json"
	"io/ioutil"
)

func Load(filename string) (*C, error) {
	var c C
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dat, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
