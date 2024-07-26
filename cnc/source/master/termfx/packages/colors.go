package packages

import (
	"fmt"
	"strconv"
	"strings"
)

type Colors struct {
	Red, DarkRed         string
	Blue, DarkBlue       string
	Green, DarkGreen     string
	Yellow, DarkYellow   string
	Magenta, DarkMagenta string
	Cyan, DarkCyan       string
	Gray, DarkGray       string
	PuttyWhite, White    string
	Reset                string

	Colorize func(hex string) string
}

func Colorize(hex string) string {
	rgb, err := Hex2RGB(hex)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("\u001B[38;2;%d;%d;%dm", rgb.Red, rgb.Green, rgb.Blue)
}

func ColorizeBG(hex string) string {
	rgb, err := Hex2RGB(hex)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("\u001B[48;2;%d;%d;%dm", rgb.Red, rgb.Green, rgb.Blue)
}

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func Hex2RGB(hex string) (RGB, error) {
	if strings.HasPrefix(string(hex), "#") {
		hex = hex[1:]
	}

	var rgb RGB
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}
