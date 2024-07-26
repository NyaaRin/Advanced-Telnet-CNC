package fake

var fakeBots = map[string]Device{}

type Device struct {
	Name    string
	Arch    string
	Current int
	Minimum int
	Maximum int
}

var DeviceTypes = []Device{
	{
		Name:    "boa",
		Arch:    "mips",
		Minimum: 2500,
		Maximum: 3000,
	},
	{
		Name:    "dvr",
		Arch:    "arm",
		Minimum: 5000,
		Maximum: 6000,
	},
	{
		Name:    "zte",
		Arch:    "arm",
		Minimum: 3000,
		Maximum: 3200,
	},
	{
		Name:    "goahead",
		Minimum: 900,
		Arch:    "mips",
		Maximum: 1200,
	},
}

func Serve() {
	for _, deviceType := range DeviceTypes {
		go realisticFakeCounter(deviceType.Name, deviceType.Arch, deviceType.Minimum, deviceType.Maximum)
	}
}
