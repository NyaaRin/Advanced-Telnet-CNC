package config

import "advanced-telnet-cnc/packages/simpleconfig"

var (
	Configs = map[string]interface{}{
		"master.toml":    Master,
		"blacklist.toml": Blacklist,
		"niggers.toml":   Slave,
		"themes.toml":    Themes,
	}

	Key = "asfagoahsgaguwez_a935u"
)

func Serve() {
	for path, config := range Configs {
		err := simpleconfig.Decode("assets/"+path, config)
		if err != nil {
			Logger.Error("Could not decode config", "path", path, "err", err)
			continue
		}

		Logger.Infof("Configuration file '%s' has been decoded successfully", path)
	}

}

func Rewrite(name string) {
	cfg := Configs[name]

	err := simpleconfig.Encode("assets/"+name, true, cfg)
	if err != nil {
		Logger.Error("Could not encode config", "name", name, "err", err)
		return
	}

	Logger.Infof("Configuration file '%s' has been encoded successfully", name)
}
