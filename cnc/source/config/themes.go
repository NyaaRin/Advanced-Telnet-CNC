package config

var (
	Themes = &themes{
		Themes: []*Theme{
			{
				Name:      "default",
				Primary:   "#ffffff",
				Secondary: "#b546e8",
			},
		},
	}
)

type themes struct {
	Themes []*Theme `toml:"themes"`
}

type Theme struct {
	Name      string `toml:"name"`
	Primary   string `toml:"primary"`
	Secondary string `toml:"secondary"`
}

func ThemeByName(name string) *Theme {
	for _, theme := range Themes.Themes {
		if theme.Name == name {
			return theme
		}
	}

	return nil
}
