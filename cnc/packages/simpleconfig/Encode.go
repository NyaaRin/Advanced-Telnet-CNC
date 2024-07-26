package simpleconfig

import (
	"github.com/pelletier/go-toml"
	"os"
)

func Encode(path string, overwrite bool, v interface{}) error {
	if !overwrite {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return err
		}
	}

	createdFile, err := os.Create(path)
	if err != nil {
		return err
	}

	err = toml.NewEncoder(createdFile).Encode(v)
	if err != nil {
		return err
	}

	return nil
}
