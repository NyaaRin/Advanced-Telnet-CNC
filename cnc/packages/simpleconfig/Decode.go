package simpleconfig

import (
	"bytes"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

func Decode(path string, v interface{}) error {
	err := Encode(path, false, v)
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.NewDecoder(bytes.NewReader(file)).Decode(v)
	if err != nil {
		return err
	}

	return nil
}
