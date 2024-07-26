package fake

func Distribution() map[string]int {
	var arches = make(map[string]int)
	for _, device := range fakeBots {
		arches[device.Name] += device.Current
	}

	return arches
}
func Arches() map[string]int {
	res := make(map[string]int)

	for _, slave := range fakeBots {
		res[slave.Arch] += slave.Current
	}
	return res
}

func Count() int {
	var count = 0
	for _, device := range fakeBots {
		count += device.Current
	}
	return count
}
