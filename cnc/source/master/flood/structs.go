package flood

import "reflect"

type Method struct {
	ID          int
	Name        string
	Aliases     []string
	Description string
	Options     []uint8
}

type AttackProfile struct {
	Id       int
	Duration int
	Targets  map[uint32]uint8
	Options  map[uint8]string
}

type Option struct {
	ID          uint8
	Description string
	Type        reflect.Kind
	Maximum     int
	Minimum     int
	Default     interface{}
	HasDefault  bool
}
