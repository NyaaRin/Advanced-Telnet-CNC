package config

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	OrigRed = "\u001B[31m"

	Red     = Colorize("#ff5555")
	DarkRed = Colorize("#bb0000")

	Green     = Colorize("#16C60C")
	DarkGreen = Colorize("#00bb00")

	Yellow     = Colorize("#ffff55")
	DarkYellow = Colorize("#bbbb00")

	Blue     = Colorize("#5555ff")
	DarkBlue = Colorize("#0000bb")

	Magenta     = Colorize("#ff55ff")
	DarkMagenta = Colorize("#CE4E9F")

	Cyan     = Colorize("#89CAFF")
	DarkCyan = Colorize("#00bbbb")

	PuttyWhite = Colorize("#bbbbbb")
	White      = Colorize("#FFFFFF")
	Gray       = Colorize("#7f7f7f")
	DarkGray   = Colorize("#474747")

	Primary   = White
	Secondary = Cyan

	Reset = "\033[0m"

	MOTD = ""
)

func Colorize(hex Hex) string {
	rgb, err := Hex2RGB(hex)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("\u001B[38;2;%d;%d;%dm", rgb.Red, rgb.Green, rgb.Blue)
}

type Hex string

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (h Hex) toRGB() (RGB, error) {
	return Hex2RGB(h)
}

func Hex2RGB(hex Hex) (RGB, error) {
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
